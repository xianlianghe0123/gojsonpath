package ast

import (
	"fmt"
	"reflect"
)

type IndexField struct {
	Index int
}

func NewIndexField(index int) *IndexField {
	return &IndexField{
		Index: index,
	}
}

func (i *IndexField) String() string {
	return fmt.Sprintf("[%d]", i.Index)
}

func (i *IndexField) SingleResult() bool {
	return true
}

func (i *IndexField) Get(data interface{}) (interface{}, error) {
	value := reflect.ValueOf(data)
	for value.Type().Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil, i.errNotFound()
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Array && value.Kind() != reflect.Slice {
		return nil, i.errNotFound()
	}
	idx := (i.Index + value.Len()) % value.Len()
	if idx < 0 || idx >= value.Len() {
		return nil, i.errNotFound()
	}
	return value.Index(idx).Interface(), nil
}

func (i *IndexField) errNotFound() error {
	return fmt.Errorf("%d not found", i.Index)
}
