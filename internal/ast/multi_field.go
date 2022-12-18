package ast

import (
	"fmt"
	"reflect"
	"strings"
)

type MultiFields struct {
	Fields []string
}

func NewMultiFields(fields ...string) *MultiFields {
	return &MultiFields{
		Fields: fields,
	}
}

func (m *MultiFields) String() string {
	builder := strings.Builder{}
	builder.WriteRune('[')
	for i, f := range m.Fields {
		builder.WriteString(fmt.Sprintf("%q", f))
		if i < len(m.Fields)-1 {
			builder.WriteRune(',')
		}
	}
	builder.WriteRune(']')
	return builder.String()
}

func (m *MultiFields) SingleResult() bool {
	return false
}

func (m *MultiFields) Get(data interface{}) (interface{}, error) {
	value := reflect.ValueOf(data)
	for value.Type().Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil, fmt.Errorf("can not find %s from nil", m)
		}
		value = value.Elem()
	}
	var (
		result []interface{}
		err    error
	)
	switch value.Kind() {
	case reflect.Map, reflect.Struct:
		result, err = m.getObject(value)
	default:
		return nil, fmt.Errorf("unsupported get field %s from %s", m, value.Type().Kind())
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MultiFields) getObject(value reflect.Value) ([]interface{}, error) {
	data := value.Interface()
	result := make([]interface{}, 0, len(m.Fields))
	for _, field := range m.Fields {
		v, err := NewSingleField(field).Get(data)
		if err != nil {
			continue
		}
		result = append(result, v)
	}
	return result, nil
}
