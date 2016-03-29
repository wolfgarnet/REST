package REST

type Session struct {
	User User
	Data map[string]string
}

func (s *Session) Get(key, def string) string {
	val, ok := s.Data[key]
	if !ok {
		return def
	}

	return val
}
