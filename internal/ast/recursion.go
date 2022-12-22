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
	result := make([]interface{}, 0)
	result, err := r.get(reflect.ValueOf(data), result)
	if err != nil {
		return nil, err
	}
	return &Result{
		data:  result,
		multi: true,
	}, nil
}

func (r *Recursion) get(value reflect.Value, result []interface{}) ([]interface{}, error) {
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil, nil
		}
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String:
		return result, nil
	case reflect.Array, reflect.Slice:
		return r.getArray(value, result), nil
	case reflect.Map:
		return r.getMap(value, result), nil
	case reflect.Struct:
		return r.getStruct(value, result), nil
	case reflect.Interface:
		return r.get(reflect.ValueOf(value.Interface()), result)
	default:
		return nil, fmt.Errorf("unsupported get field %s from %s", r, value.Kind().String())
	}
}

func (r *Recursion) getMap(value reflect.Value, result []interface{}) []interface{} {
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
		r, err := r.get(iter.Value(), result)
		if err != nil {
			continue
		}
		result = r
	}
	return result
}

func (r *Recursion) getStruct(value reflect.Value, result []interface{}) []interface{} {
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
		r, err := r.get(value.Field(i), result)
		if err != nil {
			continue
		}
		result = r
	}
	return result
}

func (r *Recursion) getArray(value reflect.Value, result []interface{}) []interface{} {
	t, err := r.next.Get(value.Interface())
	if err == nil {
		if t.multi {
			result = append(result, t.data.([]interface{})...)
		} else {
			result = append(result, t.data)
		}
	}
	for i := 0; i < value.Len(); i++ {
		r, err := r.get(value.Index(i), result)
		if err != nil {
			continue
		}
		result = r
	}
	return result
}
