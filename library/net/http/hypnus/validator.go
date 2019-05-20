package hypnus

import (
	"reflect"
	"sync"

	"gopkg.in/go-playground/validator.v9"
)

type defaultValidator struct {
	once     sync.Once
	Validate *validator.Validate
}

func (v *defaultValidator) lazyInit() {
	v.once.Do(func() {
		v.Validate = validator.New()
	})
}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	v.lazyInit()
	if isKindOfStruct(obj) {
		if err := v.Validate.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}

func isKindOfStruct(obj interface{}) bool {

	typeof := reflect.TypeOf(obj)

	if typeof.Kind() == reflect.Struct {
		return true
	}
	if typeof.Kind() == reflect.Ptr && typeof.Elem().Kind() == reflect.Struct {
		return true
	}
	return false
}
