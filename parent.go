package REST

// Parent represents a node that has a child url node
type Parent interface {
	// GetChild gets the child of the parent given its name
	GetChild(name string) (Node, error)
}