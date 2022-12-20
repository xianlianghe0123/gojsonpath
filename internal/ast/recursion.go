package ast

import (
	"fmt"
	"reflect"
)

type Recursion struct {
	next Node
}

func NewRecursion(next Node) *Recursion {
	return &Recursion{
		next: next,
	}
}

func (r *Recursion) String() string {
	return fmt.Sprintf("..%s", r.next.String())
}

func (r *Recursion) Get(data interface{}) (*Result, error) {
	re, err := r.get(reflect.ValueOf(data))
	if err != nil {
		return nil, err
	}
	return &Result{
		data:  re,
		multi: true,
	}, nil
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
		result, err := r.next.Get(value.Interface())
		if err != nil {
			return []interface{}{}, nil
		}
		if result.multi {
			return result.data.([]interface{}), nil
		} else {
			return []interface{}{result.data}, nil
		}
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
	t, err := r.next.Get(value.Interface())
	if err == nil {
		if t.multi {
			result = append(result, t.data.([]interface{})...)
		} else {
			result = append(result, t.data)
		}
	}
	iter := value.MapRange()
	for iter.Next() {
		e, _ := r.get(iter.Value())
		result = append(result, e...)
	}
	return result
}

func (r *Recursion) getStruct(value reflect.Value) []interface{} {
	result := make([]interface{}, 0, value.NumField()+1)
	t, err := r.next.Get(value.Interface())
	if err == nil {
		if t.multi {
			result = append(result, t.data.([]interface{})...)
		} else {
			result = append(result, t.data)
		}
	}
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
	t, err := r.next.Get(value.Interface())
	if err == nil {
		if t.multi {
			result = append(result, t.data.([]interface{})...)
		} else {
			result = append(result, t.data)
		}
	}
	for i := 0; i < value.Len(); i++ {
		e, _ := r.get(value.Index(i))
		result = append(result, e...)

	}
	return result
}
