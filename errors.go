package di

import (
	"fmt"
	"reflect"
)

type NotAFunc struct {
}

func (e NotAFunc) Error() string {
	return "source is not a function"
}

type DependencyNotFound struct {
	requested reflect.Type
}

func (e DependencyNotFound) Error() string {
	return fmt.Sprintf("dependency %s not found", e.requested.String())
}

type CircularDependency struct {
}

func (e CircularDependency) Error() string {
	return "circular dependency found"
}

type InvalidInvocation struct {
}

func (e InvalidInvocation) Error() string {
	return "cannot invoke outside init scope"
}
