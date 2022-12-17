package ast

import "strings"

type IndexesField struct {
	nodes []Node
}

func NewIndexesField(nodes ...Node) *IndexesField {
	return &IndexesField{
		nodes: nodes,
	}
}

func (i *IndexesField) String() string {
	builder := strings.Builder{}
	builder.WriteRune('[')
	for j, node := range i.nodes {
		t := node.String()
		t = t[1 : len(t)-1]
		builder.WriteString(t)
		if j < len(i.nodes)-1 {
			builder.WriteRune(',')
		}
	}
	builder.WriteRune(']')
	return builder.String()
}

func (i *IndexesField) SingleResult() bool {
	return false
}

func (i *IndexesField) Get(data interface{}) (interface{}, error) {
	result := make([]interface{}, 0)
	for _, n := range i.nodes {
		r, err := n.Get(data)
		if err != nil {
			continue
		}
		if n.SingleResult() {
			result = append(result, r)
		} else {
			result = append(result, r.([]interface{})...)
		}

	}
	return result, nil
}
