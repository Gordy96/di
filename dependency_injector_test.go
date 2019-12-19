package di

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type Foo struct {
	i string
}

func (f *Foo) Do() {
	fmt.Printf("now is: %s\n", f.i)
}

func TestDI(t *testing.T) {
	var err error
	di := DI{map[reflect.Kind]*entry{}}
	err = di.Register(func(i string) *Foo {
		return &Foo{i}
	})
	if err != nil {
		t.Fail()
	}
	err = di.Register(func () string {
		return time.Now().String()
	})
	if err != nil {
		t.Fail()
	}

	err = di.Ensure()
	if err != nil {
		t.Fail()
	}

	err = di.Get(func(f *Foo) {
		f.Do()
	})
	if err != nil {
		t.Fail()
	}
}
