package REST

type Extension interface {
	GetName() string
}

type Autonomous interface  {
	Autonomize(context *Context) (Response, error)
}

type Node interface {
	Identifier() string
	UrlName() string
	Parent() Node
}

type Savable interface  {
	Save(*Context)
}

type Resource interface {
	Node
	Savable

	GetDisplayName() string
	GetExtensions() []Extension

	// Type returns a lower case name of the type. Eg. user, file etc
	Type() string
	//Type() reflect.Type

	/*
	setId(t string)
	setType(t string)
	*/
	//SetIdentifier(id string)

	Describe() string
}

type Action interface {
	Extension

	IsApplicable(node Node) bool
}