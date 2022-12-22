package ast

import (
	"fmt"
	"reflect"
	"strconv"
)

type SingleField struct {
	field string
	next  Node
}

func NewSingleField(field string, next Node) *SingleField {
	return &SingleField{
		field: field,
		next:  next,
	}
}

func (s *SingleField) String() string {
	return fmt.Sprintf("[%q]%s", s.field, s.next.String())
}

func (s *SingleField) Get(data interface{}) (*Result, error) {
	value := reflect.ValueOf(data)
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil, s.errNotFound()
		}
		value = value.Elem()
	}
	switch value.Kind() {
	case reflect.Map:
		return s.getMap(value)
	case reflect.Struct:
		return s.getStruct(value)
	default:
		return nil, fmt.Errorf("unsupported get field %s from %s", s.field, value.Kind())
	}
}

func (s *SingleField) errNotFound() error {
	return fmt.Errorf("%s not found", s.field)
}

func (s *SingleField) getMap(value reflect.Value) (*Result, error) {
	key := reflect.ValueOf(s.field)
	switch t := value.Type().Key().Kind(); t {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s.field, 10, 64)
		if err != nil {
			return nil, s.errNotFound()
		}
		switch t {
		case reflect.Int:
			key = reflect.ValueOf(int(i))
		case reflect.Int8:
			key = reflect.ValueOf(int8(i))
		case reflect.Int16:
			key = reflect.ValueOf(int16(i))
		case reflect.Int32:
			key = reflect.ValueOf(int32(i))
		case reflect.Int64:
			key = reflect.ValueOf(i)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(s.field, 10, 64)
		if err != nil {
			return nil, s.errNotFound()
		}
		switch t {
		case reflect.Uint:
			key = reflect.ValueOf(uint(u))
		case reflect.Uint8:
			key = reflect.ValueOf(uint8(u))
		case reflect.Uint16:
			key = reflect.ValueOf(uint16(u))
		case reflect.Uint32:
			key = reflect.ValueOf(uint32(u))
		case reflect.Uint64:
			key = reflect.ValueOf(u)
		}
	case reflect.Float32, reflect.Float64:
		float, err := strconv.ParseFloat(s.field, 64)
		if err != nil {
			return nil, s.errNotFound()
		}
		switch t {
		case reflect.Float32:
			key = reflect.ValueOf(float32(float))
		case reflect.Float64:
			key = reflect.ValueOf(float)
		}
	case reflect.Bool:
		b, err := strconv.ParseBool(s.field)
		if err != nil {
			return nil, s.errNotFound()
		}
		key = reflect.ValueOf(b)
	case reflect.String:
	default:
		return nil, fmt.Errorf("unsupported map where key type is %s", t)
	}
	v := value.MapIndex(key)
	if !v.IsValid() {
		return nil, s.errNotFound()
	}
	return s.next.Get(v.Interface())
}

func (s *SingleField) getStruct(value reflect.Value) (*Result, error) {
	for i := 0; i < value.NumField(); i++ {
		key, omitempty := getFieldKey(value.Type().Field(i))
		if key != s.field {
			continue
		}
		if omitempty && value.IsZero() {
			break
		}
		return s.next.Get(value.Field(i).Interface())
	}
	return nil, s.errNotFound()
}
