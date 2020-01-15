package di

import (
	mapset "github.com/deckarep/golang-set"
	"reflect"
)

type entry struct {
	constructor interface{}
	instance    interface{}
	singleton   bool
}

func NewDI() *DI {
	return &DI{
		entries: map[reflect.Type]*entry{},
		_seal:   true,
	}
}

type DI struct {
	entries map[reflect.Type]*entry
	_seal   bool
}

func (di *DI) Singleton(source interface{}) error {
	return di.put(source, true)
}

func (di *DI) Factory(source interface{}) error {
	return di.put(source, false)
}

func (di *DI) Ensure() error {
	return di.resolve()
}

func (di *DI) Register(callback func(*DI) error) error {
	di.unseal()
	defer di.seal()
	err := callback(di)
	if err != nil {
		return err
	}
	err = di.resolve()
	return err
}

func (di *DI) Get(callback interface{}) error {
	cbType := reflect.TypeOf(callback)
	cbKind := cbType.Kind()
	if cbKind != reflect.Func {
		return NotAFunc{}
	}
	requestedType := cbType.In(0)

	entry, found := di.entries[requestedType]
	if !found {
		return DependencyNotFound{requestedType}
	}

	temp, err := di.create(entry)
	if err != nil {
		return err
	}
	reflect.ValueOf(callback).Call([]reflect.Value{reflect.ValueOf(temp)})
	return nil
}

func (di *DI) sealed() bool {
	return di._seal
}
func (di *DI) seal() {
	di._seal = true
}
func (di *DI) unseal() {
	di._seal = false
}

func (di *DI) put(source interface{}, singleton bool) error {
	if di.sealed() {
		return InvalidInvocation{}
	}
	sourceType := reflect.TypeOf(source)
	sourceKind := sourceType.Kind()
	if sourceKind != reflect.Func {
		return NotAFunc{}
	}
	retType := sourceType.Out(0)
	di.entries[retType] = &entry{source, nil, singleton}
	return nil
}

func (di *DI) resolve() error {
	graph := make(map[reflect.Type]mapset.Set)
	for entryType, entry := range di.entries {
		dependencySet := mapset.NewSet()
		constructorType := reflect.TypeOf(entry.constructor)
		for i := 0; i < constructorType.NumIn(); i++ {
			dependencySet.Add(constructorType.In(i))
		}
		graph[entryType] = dependencySet
	}
	for len(graph) != 0 {
		readySet := mapset.NewSet()
		for entryType, deps := range graph {
			if deps.Cardinality() == 0 {
				readySet.Add(entryType)
			}
		}
		if readySet.Cardinality() == 0 {
			return CircularDependency{}
		}
		for entryType := range readySet.Iter() {
			delete(graph, entryType.(reflect.Type))
		}
		for entryType, deps := range graph {
			diff := deps.Difference(readySet)
			graph[entryType] = diff
		}
	}

	return nil
}

func (di *DI) create(entry *entry) (interface{}, error) {
	if entry.singleton && entry.instance != nil {
		return entry.instance, nil
	}
	var deps []reflect.Value
	constructorType := reflect.TypeOf(entry.constructor)
	for i := 0; i < constructorType.NumIn(); i++ {
		inType := constructorType.In(i)
		in, found := di.entries[inType]
		if !found {
			return nil, DependencyNotFound{inType}
		}
		temp, err := di.create(in)
		if err != nil {
			return nil, err
		}
		deps = append(deps, reflect.ValueOf(temp))
	}
	ret := reflect.ValueOf(entry.constructor).Call(deps)
	if len(ret) > 1 {
		return nil, ret[1].Interface().(error)
	}
	if ret[0].IsValid() {
		if entry.singleton {
			entry.instance = ret[0].Interface()
		}
		return ret[0].Interface(), nil
	}
	return nil, DependencyNotFound{constructorType.Out(0)}
}
