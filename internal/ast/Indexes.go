package ast

import (
	"strings"
)

type Indexes struct {
	nodes []Node
}

func NewIndexes(nodes []Node) *Indexes {
	return &Indexes{
		nodes: nodes,
	}
}

func (i *Indexes) String() string {
	builder := strings.Builder{}
	builder.WriteRune('[')
	for j, node := range i.nodes {
		builder.WriteString(strings.SplitN(node.String(), "]", 2)[0][1:])
		if j < len(i.nodes)-1 {
			builder.WriteRune(',')
		}
	}
	builder.WriteRune(']')
	return builder.String()
}

func (i *Indexes) Get(data interface{}) (*Result, error) {
	r, err := i.get(data)
	if err != nil {
		return nil, err
	}
	return &Result{
		data:  r,
		multi: true,
	}, nil
}

func (i *Indexes) get(data interface{}) ([]interface{}, error) {
	result := make([]interface{}, 0)
	for _, n := range i.nodes {
		r, err := n.Get(data)
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
