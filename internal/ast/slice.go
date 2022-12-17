package ast

import (
	"fmt"
	"reflect"
	"strconv"
)

type SliceField struct {
	Start int
	End   int
	Step  int
}

func NewSliceField(start, end, step int) *SliceField {
	return &SliceField{
		Start: start,
		End:   end,
		Step:  step,
	}
}

func (s *SliceField) String() string {
	start, end, step := "", "", ""
	if s.Start != 0 {
		start = strconv.Itoa(s.Start)
	}
	if s.End != 0 {
		end = strconv.Itoa(s.End)
	}
	if s.Step != 0 {
		step = strconv.Itoa(s.Step)
	}
	return fmt.Sprintf("[%s:%s:%s]", start, end, step)
}

func (s *SliceField) SingleResult() bool {
	return false
}

func (s *SliceField) Get(data interface{}) (interface{}, error) {
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Array && value.Kind() != reflect.Slice {
		return nil, fmt.Errorf("can not get slice without array")
	}
	start := (s.Start + value.Len()) % value.Len()
	if start < 0 {
		start = 0
	}
	end := (s.End + value.Len()) % value.Len()
	if end < 0 {
		return nil, fmt.Errorf("end index out of bounds: [%d]", s.End)
	} else if end == 0 {
		end = value.Len()
	}
	if start >= end {
		return nil, fmt.Errorf("not found")
	}
	step := s.Step
	if step == 0 {
		step = 1
	}
	result := make([]interface{}, 0, (end-start)/step)
	for i := start; i < value.Len() && i < end; i += step {
		result = append(result, value.Index(i).Interface())
	}
	return result, nil
}
