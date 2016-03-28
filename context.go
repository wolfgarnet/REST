package REST

import (
	"net/http"
	"bytes"
	"strconv"
)

type Context struct {
	Request *http.Request
	Session *Session
	System System
}

func newContext(request *http.Request, system *System) *Context {
	session := newSession(request)
	return &Context{request, session, system}
}

func (c *Context) GetInt(field string, value int) int {
	v, err := strconv.Atoi(c.Request.FormValue(field))
	if err != nil {
		return value
	}

	return v
}