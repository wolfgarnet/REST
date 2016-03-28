package REST

// Super is a pseudo inheritance interface, returning the "super" type of the
// current implementation.
type Super interface {
	Super() interface{}
}