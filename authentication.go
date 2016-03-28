package REST

import (
	"net/http"
)

// Authentication represents ...
type Authentication interface {
	Authenticate(request *http.Request) User
	Login(username, password string)
}

type Permission uint8;

const (
	AUTH_NONE Permission = iota
)

// ACL represents an access control list
type ACL interface {
	GetPermission(user User) Permission;
}