package REST

import (
	"net/http"
)

// Authentication represents ...
type Authentication interface {
	Authenticate(request *http.Request) User
	Login(username, password string)
}