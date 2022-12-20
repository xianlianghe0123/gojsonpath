package ast

type End struct{}

func NewEnd() End {
	return End{}
}

func (e End) String() string {
	return ""
}

func (e End) Get(data interface{}) (*Result, error) {
	return &Result{
		data:  data,
		multi: false,
	}, nil
}
