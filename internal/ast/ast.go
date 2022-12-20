package ast

type Result struct {
	data  interface{}
	multi bool
}

type Node interface {
	Get(interface{}) (*Result, error)
	String() string
}

type AST struct {
	node Node
}

func NewAST(node Node) *AST {
	return &AST{
		node: node,
	}
}

func (a *AST) String() string {
	if a.node == nil {
		return ""
	}
	return a.node.String()
}

func (a *AST) Get(data interface{}) (interface{}, error) {
	if a.node == nil {
		return data, nil
	}
	result, err := a.node.Get(data)
	if err != nil {
		return nil, err
	}
	return result.data, nil
}
