package REST

import (
	"web/core/response"
	"github.com/wolfgarnet/templates"
	"strings"
)


type UrlMethod func(context *Context) (response.Response, error)

type UrlMethodRetriever interface {
	GetUrlMethod(methodName, method string) UrlMethod
}

type Runner interface {
	Run(context *Context) (response.Response, error)
}

func GetMethodRunner(object interface{}, urlName, method string) Runner {
	logger.Debug("GET METHOD RUNNER: %v for %v", urlName, object)
	um, ok := object.(UrlMethodRetriever)
	if !ok {
		return nil
	}

	f := um.GetUrlMethod(strings.ToLower(urlName), strings.ToLower(method))
	if f == nil {
		return nil
	}

	return &MethodRunner{f}
}

// REDIRECT

type Redirect struct  {
	Url string
}

func (runner *Redirect) Run(context *Context) (response.Response, error) {
	return response.NewRedirectResponse(runner.Url), nil
}

// TEMPLATE

type Template struct {
	Renderer *templates.Renderer
}

func (runner *Template) Run(context *Context) (response.Response, error) {
	r := response.NewBufferedResponse()
	context.Request.ParseForm()

	runner.Renderer.AddData("system", context.System)
	runner.Renderer.AddData("context", context)
	err := runner.Renderer.Render(r.Bytes)
	return r, err
}

// METHOD

type MethodRunner struct {
	f UrlMethod
}

func (runner *MethodRunner) Run(context *Context) (response.Response, error) {
	logger.Debug("RUNNING METHOD RUNNER ==== %v", runner.f)
	return runner.f(context)
}


/// AUTONOMOUS

type AutonomousRunner struct {
	Node Autonomous
	Context *Context
}

func (runner *AutonomousRunner) Run(context *Context) (response.Response, error) {
	println("RUNNING AUTONOMOUS RUNNER")
	return runner.Node.Autonomize(runner.Context)
}

/// ERROR

func MakeErrorResponse(error string) response.Response {
	b := response.NewBufferedResponse()

	b.Status = 500
	b.Append(error)

	return b
}



type ErrorRunner struct {
	Status int
	Message string
}

/*
func NewErrorRunner(status int, message string) response.Response {
	return &ErrorRunner{status, message}
}
*/

func (runner *ErrorRunner) Run(context *Context) (response.Response, error) {
	r := response.NewBufferedResponse()
	r.Append(runner.Message)
	r.SetStatus(runner.Status)
	r.Status = runner.Status
	logger.Debug("I WAS ERE: %v", r)
	return r, nil
}

