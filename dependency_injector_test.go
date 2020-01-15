package di

import (
	"fmt"
	"testing"
	"time"
)

type Foo struct {
	i string
}

func (f *Foo) Do() {
	fmt.Printf("now is: %s\n", f.i)
}

type Bar struct {
}

func (b *Bar) Do() {

}

type Doer interface {
	Do()
}

func TestDI(t *testing.T) {
	var err error
	di := NewDI()
	err = di.Register(func(di *DI) error {
		err = di.Factory(func(i string) Doer {
			return &Foo{i}
		})
		if err != nil {
			return err
		}
		err = di.Factory(func() string {
			return time.Now().String()
		})
		return err
	})
	if err != nil {
		t.Fail()
	}

	err = di.Get(func(f Doer) {
		f.Do()
	})
	if err != nil {
		t.Fail()
	}
	err = di.Get(func(b *Bar) {
		b.Do()
	})
	if err == nil {
		t.Fail()
	} else {
		t.Log(err)
	}
}
