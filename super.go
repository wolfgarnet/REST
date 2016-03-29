package REST

import "reflect"

// Super is a pseudo inheritance interface, returning the "super" type of the
// current implementation.
type Super interface {
	Super() interface{}
}

func GetSuperType(object interface{}) interface{} {
	t := reflect.TypeOf(object)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	superType := reflect.TypeOf((*Super)(nil)).Elem()
	if t.Implements(superType) {
		logger.Debug("Supertype: %v", superType)
		sm, ok := object.(Super)
		logger.Debug("SM: %v -- %v", sm, ok)
		if !ok {
			return nil
		}

		return sm
	}

	return nil
}