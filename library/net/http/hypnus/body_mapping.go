package hypnus

import (
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var scache = &cache{
	data: make(map[reflect.Type]*sinfo),
}

type sinfo struct {
	field []*field
}

type cache struct {
	data  map[reflect.Type]*sinfo
	mutex sync.RWMutex
}

type field struct {
	StructField  reflect.StructField
	name         string
	hasDefault   bool
	defaultValue reflect.Value
}

func mapBody(ptr interface{}, body map[string]string) error {
	sinfo := scache.get(reflect.TypeOf(ptr))
	val := reflect.ValueOf(ptr).Elem()

	for i, field := range sinfo.field {
		typeField := field.StructField
		structField := val.Field(i)
		if !structField.CanSet() {
			continue
		}

		typeFieldKind := typeField.Type.Kind()
		inputValue, exist := body[field.name]

		// don't handle struct type now . just handle simple type: string 、 number 、 bool、time
		if typeFieldKind == reflect.Struct {

			/*
				var t time.Time = time.Time{}
				reflect.TypeOf(t).Kind() == reflect.Struct
				so specifically to handle time.Time
			*/
			if _, isTime := structField.Interface().(time.Time); isTime && exist {
				if err := setTimeValue(inputValue, structField); err != nil {
					return err
				}
			}
			continue
		}

		if !exist {
			if field.hasDefault {
				structField.Set(field.defaultValue)
				continue
			}
		}

		// if inputValue is empty, set field with defaultValue
		if field.hasDefault && inputValue == "" {
			structField.Set(field.defaultValue)
			continue
		}

		if err := setFieldValue(typeFieldKind, inputValue, structField); err != nil {
			return err
		}
	}
	return nil
}

// get get sinfo from cache.
func (c *cache) get(obj reflect.Type) (s *sinfo) {
	var ok bool
	c.mutex.RLock()
	if s, ok = c.data[obj]; !ok {
		c.mutex.RUnlock()
		s = c.set(obj)
		return
	}
	c.mutex.RUnlock()
	return
}

func (c *cache) set(obj reflect.Type) (s *sinfo) {
	s = new(sinfo)
	tp := obj.Elem()
	for i := 0; i < tp.NumField(); i++ {
		fd := new(field)
		fd.StructField = tp.Field(i)
		tag := fd.StructField.Tag.Get("body")
		if tag == "" {
			fd.name = fd.StructField.Name
		} else {
			fd.name = tag
		}
		s.field = append(s.field, fd)

		if defV := fd.StructField.Tag.Get("default"); defV != "" {
			fd.hasDefault = true

			// NOTE : don't call of reflect.Value.Elem on zero Value
			// reflect.New(reflect.Type) return pointer of the new value
			val := reflect.New(fd.StructField.Type).Elem()
			setFieldValue(fd.StructField.Type.Kind(), defV, val)
			fd.defaultValue = val
		}
	}
	c.mutex.Lock()
	c.data[obj] = s
	c.mutex.Unlock()
	return
}

func setFieldValue(fieldKind reflect.Kind, inputValue string, val reflect.Value) error {
	switch fieldKind {
	case reflect.Int:
		return setIntValue(inputValue, 0, val)
	case reflect.Int8:
		return setIntValue(inputValue, 8, val)
	case reflect.Int16:
		return setIntValue(inputValue, 16, val)
	case reflect.Int32:
		return setIntValue(inputValue, 32, val)
	case reflect.Int64:
		return setIntValue(inputValue, 64, val)
	case reflect.Bool:
		return setBoolValue(inputValue, val)
	case reflect.String:
		return setStringValue(inputValue, val)
	case reflect.Float32:
		return setFloatValue(inputValue, 32, val)
	case reflect.Float64:
		return setFloatValue(inputValue, 64, val)
	}
	return nil
}

func setIntValue(inputValue string, baseSize int, val reflect.Value) error {
	if inputValue == "" {
		inputValue = "0"
	}
	i, err := strconv.ParseInt(inputValue, 10, baseSize)
	if err == nil {
		val.SetInt(i)
	}
	return errors.WithStack(err)
}

func setBoolValue(inputValue string, val reflect.Value) error {
	if inputValue == "" {
		inputValue = "false"
	}
	b, err := strconv.ParseBool(inputValue)
	if err == nil {
		val.SetBool(b)
	}
	return errors.WithStack(err)
}

func setFloatValue(inputValue string, baseSize int, val reflect.Value) error {
	if inputValue == "" {
		inputValue = "0"
	}
	f, err := strconv.ParseFloat(inputValue, baseSize)
	if err == nil {
		val.SetFloat(f)
	}
	return errors.WithStack(err)
}

func setStringValue(inputValue string, val reflect.Value) error {
	val.SetString(inputValue)
	return nil
}

func setTimeValue(inputValue string, val reflect.Value) error {

	// var t time.Time / var t time.Time = time.Time{}  both of they are zore value.
	// so when inputValue == " don't need to set time.Time{}

	if inputValue == "" {
		return nil
	}

	// why the date format is 2006-01-02 15:04:05, see https://www.jianshu.com/p/c7f7fbb16932
	t, err := time.ParseInLocation("2006-01-02 15:04:05", inputValue, time.Local)
	if err != nil {
		return errors.WithStack(err)
	}

	val.Set(reflect.ValueOf(t))
	return nil
}
