package ast

import (
	"fmt"
	"reflect"
)

type All struct{}

func NewAll() *All {
	return &All{}
}

func (a *All) String() string {
	return "[*]"
}

func (a *All) SingleResult() bool {
	return false
}

func (a *All) Get(data interface{}) (interface{}, error) {
	value := reflect.ValueOf(data)
	for value.Type().Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil, fmt.Errorf("can not find * from nil")
		}
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Map:
		return a.getMap(value)
	case reflect.Struct:
		return a.getStruct(value)
	case reflect.Slice, reflect.Array:
		return a.getArray(value)
	default:
		return nil, fmt.Errorf("unsupported find * from %s", value.Type().Kind())
	}
}

func (a *All) getMap(value reflect.Value) ([]interface{}, error) {
	result := make([]interface{}, 0, value.Len())
	iter := value.MapRange()
	for iter.Next() {
		result = append(result, iter.Value().Interface())
	}
	return result, nil
}

func (a *All) getStruct(value reflect.Value) ([]interface{}, error) {
	result := make([]interface{}, 0, value.NumField())
	for i := 0; i < value.NumField(); i++ {
		_, omitempty := getFieldKey(value.Type().Field(i))
		if omitempty && value.Field(i).IsZero() {
			continue
		}
		result = append(result, value.Field(i).Interface())
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("empty struct")
	}
	return result, nil
}

func (a *All) getArray(value reflect.Value) ([]interface{}, error) {
	if value.Len() == 0 {
		return nil, fmt.Errorf("empty array")
	}
	result := make([]interface{}, 0, value.Len())
	for i := 0; i < value.Len(); i++ {
		result = append(result, value.Index(i).Interface())
	}
	return result, nil
}
