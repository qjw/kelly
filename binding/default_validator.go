// Copyright 2017 King Qiu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
// https://github.com/qjw/kelly

package binding

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
	"sync"
)

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ StructValidator = &defaultValidator{}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyinit()
		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}
	return nil
}

/*
func isAlpha(ctx context.Context, fl validator.FieldLevel) bool {
	alphaRegex := regexp.MustCompile("^[0-9]{4}-[0-9]{2}-[0-9]{2}$")
	return alphaRegex.MatchString(fl.Field().String())
}
*/

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		// config := &validator.Config{TagName: "binding"}
		// v.validate = validator.New(config)
		v.validate = validator.New()
		v.validate.SetTagName("binding")
		// v.validate.RegisterValidationCtx("date", isAlpha)
	})
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
