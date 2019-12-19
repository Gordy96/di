package di

import "reflect"

type entry struct {
	creator     interface{}
	creatorType reflect.Type
	instance    interface{}
	singleton   bool
}

type DI struct {
	entries map[reflect.Type]*entry
}

func (di *DI) put(source interface{}, singleton bool) error {

	sourceType := reflect.TypeOf(source)
	sourceKind := sourceType.Kind()

	if sourceKind != reflect.Func {
		return NotAFunc{}
	}

	retType := sourceType.Out(0)
	//retKind := retType.Kind()

	di.entries[retType] = &entry{source, sourceType, nil, singleton}

	return nil
}

func (di *DI) Singleton(source interface{}) error {
	return di.put(source, true)
}

func (di *DI) Register(source interface{}) error {
	return di.put(source, false)
}

func (di *DI) Get(callback interface{}) error {
	cbType := reflect.TypeOf(callback)
	cbKind := cbType.Kind()

	if cbKind != reflect.Func {
		return NotAFunc{}
	}

	requestedType := cbType.In(0)
	//requestedKind := requestedType.Kind()

	entry, found := di.entries[requestedType]
	if !found {
		return DependencyNotFound{}
	}

	temp, err := di.create(entry)
	if err != nil {
		return err
	}
	reflect.ValueOf(callback).Call([]reflect.Value{reflect.ValueOf(temp)})

	return nil
}

func (di *DI) create(entry *entry) (interface{}, error) {
	if entry.singleton && entry.instance != nil {
		return entry.instance, nil
	}
	var deps []reflect.Value
	numins := entry.creatorType.NumIn()
	for i := 0; i < numins; i++ {
		inType := entry.creatorType.In(i)
		//inKind := inType.Kind()
		in, found := di.entries[inType]
		if !found {
			return nil, DependencyNotFound{}
		}
		temp, err := di.create(in)
		if err != nil {
			return nil, err
		}
		deps = append(deps, reflect.ValueOf(temp))
	}
	ret := reflect.ValueOf(entry.creator).Call(deps)
	if len(ret) > 1 {
		return nil, ret[1].Interface().(error)
	}
	if ret[0].IsValid() {
		if entry.singleton {
			entry.instance = ret[0].Interface()
		}
		return ret[0].Interface(), nil
	}
	return nil, DependencyNotFound{}
}

func (di *DI) Ensure() error {
	for _, e := range di.entries {
		_, err := di.create(e)
		if err != nil {
			return err
		}
	}
	return nil
}
