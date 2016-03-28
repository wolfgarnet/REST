package REST

import (
	"reflect"
	"log"
)

func GetSuperType(object interface{}) interface{} {

	t := reflect.TypeOf(object)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	superType := reflect.TypeOf((*Super)(nil)).Elem()
	if t.Implements(superType) {
		log.Printf("Supertype: %v", superType)
		sm, ok := object.(Super)
		log.Printf("SM: %v -- %v", sm, ok)
		if !ok {
			return nil
		}

		return sm
	}

	return nil
}
