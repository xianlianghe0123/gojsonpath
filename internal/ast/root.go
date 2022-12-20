package ast

import "fmt"

type Root struct {
	next Node
}

func NewRoot(next Node) *Root {
	return &Root{
		next: next,
	}
}

func (r *Root) String() string {
	return fmt.Sprintf("$%s", r.next.String())
}

func (r *Root) Get(data interface{}) (*Result, error) {
	return r.next.Get(data)
}
