package ast

type Root struct{}

func NewRoot() *Root {
	return &Root{}
}

func (r *Root) String() string {
	return "$"
}

func (r *Root) SingleResult() bool {
	return true
}

func (r *Root) Get(data interface{}) (interface{}, error) {
	return data, nil
}
