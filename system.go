package REST

// System must implement system specific functionality
type System interface {
	GetRunner(session *Session, object interface{}, urlName, method string) Runner
	GetExtensions(fieldName string) []interface{}

	// RenderObject is a convenience function for rendering an object
	RenderObject(object interface{}, context *Context) string

	Save(r Resource) error

	GetResource(id string) (Resource, error)

	// GetId will provide the next consecutive id for a given type string
	GetId(t string) string
}
