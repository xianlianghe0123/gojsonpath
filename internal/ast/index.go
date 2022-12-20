package ast

import (
	"fmt"
	"reflect"
)

type Index struct {
	index int
	next  Node
}

func NewIndexField(index int, next Node) *Index {
	return &Index{
		index: index,
		next:  next,
	}
}

func (i *Index) String() string {
	return fmt.Sprintf("[%d]%s", i.index, i.next.String())
}

func (i *Index) Get(data interface{}) (*Result, error) {
	value := reflect.ValueOf(data)
	for value.Type().Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil, fmt.Errorf("index %d not found", i.index)
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Array && value.Kind() != reflect.Slice {
		return nil, fmt.Errorf("could not get index %d of type %s", i.index, value.Kind())
	}
	idx := (i.index + value.Len()) % value.Len()
	if idx < 0 || idx >= value.Len() {
		return nil, fmt.Errorf("index %d not found", i.index)
	}
	return i.next.Get(value.Index(idx).Interface())
}
