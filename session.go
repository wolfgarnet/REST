package REST

import (
	"net/http"
)

type Session struct {
	User User
	Data map[string]string
}

func newSession(request *http.Request) *Session {
	sc := &Session{"default", "desktop", "default", nil}
	return sc
}