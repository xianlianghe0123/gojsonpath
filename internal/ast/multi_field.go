package ast

import (
	"fmt"
	"reflect"
	"strings"
)

type MultiFields struct {
	fields []string
	next   Node
}

func NewMultiFields(fields []string, next Node) *MultiFields {
	return &MultiFields{
		fields: fields,
		next:   next,
	}
}

func (m *MultiFields) String() string {
	builder := strings.Builder{}
	builder.WriteRune('[')
	for i, f := range m.fields {
		builder.WriteString(fmt.Sprintf("%q", f))
		if i < len(m.fields)-1 {
			builder.WriteRune(',')
		}
	}
	builder.WriteRune(']')
	builder.WriteString(m.next.String())
	return builder.String()
}

func (m *MultiFields) Get(data interface{}) (*Result, error) {
	r, err := m.get(data)
	if err != nil {
		return nil, err
	}
	return &Result{
		data:  r,
		multi: true,
	}, nil
}

func (m *MultiFields) get(data interface{}) ([]interface{}, error) {
	value := reflect.ValueOf(data)
	for value.Type().Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil, fmt.Errorf("can not find %s from nil", m)
		}
		value = value.Elem()
	}
	switch value.Kind() {
	case reflect.Map, reflect.Struct:
		return m.getObject(value)
	default:
		return nil, fmt.Errorf("unsupported get field %s from %s", m, value.Type().Kind())
	}
}

func (m *MultiFields) getObject(value reflect.Value) ([]interface{}, error) {
	data := value.Interface()
	result := make([]interface{}, 0, len(m.fields))
	for _, field := range m.fields {
		r, err := NewSingleField(field, m.next).Get(data)
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
