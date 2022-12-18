package ast

import (
	"strings"
)

type Node interface {
	Get(interface{}) (interface{}, error)
	SingleResult() bool
	String() string
}

type AST []Node

func (a AST) String() string {
	builder := strings.Builder{}
	for _, n := range a {
		builder.WriteString(n.String())
	}
	return builder.String()
}

func (a AST) Get(data interface{}) (interface{}, error) {
	var err error
	single := true
	for _, t := range a {
		if single {
			data, err = t.Get(data)
			if err != nil {
				return nil, err
			}
		} else {
			array, ok := data.([]interface{})
			if !ok {
				continue
			}
			elems := make([]interface{}, 0, len(array))
			for j := range array {
				elem, err := t.Get(array[j])
				if err != nil {
					continue
				}
				if t.SingleResult() {
					elems = append(elems, elem)
					continue
				}
				elems = append(elems, elem.([]interface{})...)
			}
			data = elems
		}
		single = single && t.SingleResult()
	}
	return data, nil
}
