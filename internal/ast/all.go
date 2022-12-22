package ast

import (
	"fmt"
	"reflect"
)

type All struct {
	next Node
}

func NewAll(next Node) *All {
	return &All{
		next: next,
	}
}

func (a *All) String() string {
	return fmt.Sprintf("[*]%s", a.next.String())
}

func (a *All) Get(r interface{}) (*Result, error) {
	r, err := a.get(r)
	if err != nil {
		return nil, err
	}
	return &Result{
		data:  r,
		multi: true,
	}, nil
}

func (a *All) get(data interface{}) ([]interface{}, error) {
	value := reflect.ValueOf(data)
	for value.Kind() == reflect.Ptr {
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
		return nil, fmt.Errorf("unsupported find * from %s", value.Kind())
	}
}

func (a *All) getMap(value reflect.Value) ([]interface{}, error) {
	result := make([]interface{}, 0, value.Len())
	iter := value.MapRange()
	for iter.Next() {
		r, err := a.next.Get(iter.Value().Interface())
		if err != nil {
			continue
		}
		if r.multi {
			result = append(result, r.data.([]interface{})...)
		} else {
			result = append(result, r.data)
		}
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
		r, err := a.next.Get(value.Field(i).Interface())
		if err != nil {
			continue
		}
		if r.multi {
			result = append(result, r.data.([]interface{})...)
		} else {
			result = append(result, r.data)
		}
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
		r, err := a.next.Get(value.Index(i).Interface())
		if err != nil {
			continue
		}
		if r.multi {
			result = append(result, r.data.([]interface{})...)
		} else {
			result = append(result, r.data)
		}
	}
	return result, nil
}
