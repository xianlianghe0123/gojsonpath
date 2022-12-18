package ast

import (
	"fmt"
	"reflect"
)

type Recursion struct{}

func NewRecursion() *Recursion {
	return &Recursion{}
}

func (r *Recursion) String() string {
	return ".."
}

func (r *Recursion) SingleResult() bool {
	return false
}

func (r *Recursion) Get(data interface{}) (interface{}, error) {
	return r.get(reflect.ValueOf(data))
}

func (r *Recursion) get(value reflect.Value) ([]interface{}, error) {
	for value.Type().Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil, nil
		}
		value = value.Elem()
	}

	switch value.Type().Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String:
		return []interface{}{value.Interface()}, nil
	case reflect.Array, reflect.Slice:
		return r.getArray(value), nil
	case reflect.Map:
		return r.getMap(value), nil
	case reflect.Struct:
		return r.getStruct(value), nil
	case reflect.Interface:
		return r.get(reflect.ValueOf(value.Interface()))
	default:
		return nil, fmt.Errorf("unsupported get field %s from %s", r, value.Type().Kind().String())
	}
}

func (r *Recursion) getMap(value reflect.Value) []interface{} {
	result := make([]interface{}, 0, value.Len()+1)
	result = append(result, value.Interface())
	iter := value.MapRange()
	for iter.Next() {
		e, _ := r.get(iter.Value())
		result = append(result, e...)
	}
	return result
}

func (r *Recursion) getStruct(value reflect.Value) []interface{} {
	result := make([]interface{}, 0, value.NumField()+1)
	result = append(result, value.Interface())
	for i := 0; i < value.NumField(); i++ {
		_, omitempty := getFieldKey(value.Type().Field(i))
		if omitempty && value.Field(i).IsZero() {
			continue
		}
		e, err := r.get(value.Field(i))
		if err != nil {
			continue
		}
		result = append(result, e...)
	}
	return result
}

func (r *Recursion) getArray(value reflect.Value) []interface{} {
	result := make([]interface{}, 0, value.Len()+1)
	result = append(result, value.Interface())
	for i := 0; i < value.Len(); i++ {
		e, _ := r.get(value.Index(i))
		result = append(result, e...)

	}
	return result
}
