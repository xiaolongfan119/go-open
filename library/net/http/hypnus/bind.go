package hypnus

import (
	"github.com/xiaolongfan119/go-open/library/ecode"
)

var Validator = &defaultValidator{}

func Bind(body map[string]string, obj interface{}) error {

	var err error
	if err = mapBody(obj, body); err != nil {
		return ecode.ParamsFormatError_1
	}
	if err = Validator.ValidateStruct(obj); err != nil {
		return ecode.ParamsFormatError_2
	}

	return nil
}
