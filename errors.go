package di

type NotAFunc struct {

}

func (e NotAFunc) Error() string {
	return "source is not a function"
}

type DependencyNotFound struct {

}

func (e DependencyNotFound) Error() string {
	return "dependency not found"
}
