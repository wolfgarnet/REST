package REST

// System must implement system specific functionality
type System interface {
	GetRunner(session *Session, object interface{}, urlName, method string) Runner
	GetExtensions(fieldName string) []interface{}

	// RenderObject is a convenience function for rendering an object
	RenderObject(object interface{}, context *Context, session *Session)
}
