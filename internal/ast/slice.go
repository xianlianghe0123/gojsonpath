package ast

import (
	"fmt"
	"reflect"
	"strconv"
)

type Slice struct {
	start int
	end   int
	step  int
	next  Node
}

func NewSlice(start, end, step int, next Node) *Slice {
	return &Slice{
		start: start,
		end:   end,
		step:  step,
		next:  next,
	}
}

func (s *Slice) String() string {
	start, end, step := "", "", ""
	if s.start != 0 {
		start = strconv.Itoa(s.start)
	}
	if s.end != 0 {
		end = strconv.Itoa(s.end)
	}
	if s.step != 0 {
		step = strconv.Itoa(s.step)
	}
	return fmt.Sprintf("[%s:%s:%s]%s", start, end, step, s.next.String())
}

func (s *Slice) Get(data interface{}) (*Result, error) {
	r, err := s.get(data)
	if err != nil {
		return nil, err
	}
	return &Result{
		data:  r,
		multi: true,
	}, nil
}

func (s *Slice) get(data interface{}) ([]interface{}, error) {
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Array && value.Kind() != reflect.Slice {
		return nil, fmt.Errorf("can not get slice without array")
	}
	start := (s.start + value.Len()) % value.Len()
	if start < 0 {
		start = 0
	}
	end := (s.end + value.Len()) % value.Len()
	if end < 0 {
		return nil, fmt.Errorf("end index out of bounds: [%d]", s.end)
	} else if end == 0 {
		end = value.Len()
	}
	if start >= end {
		return nil, fmt.Errorf("not found")
	}
	step := s.step
	if step == 0 {
		step = 1
	}
	result := make([]interface{}, 0, (end-start)/step)
	for i := start; i < value.Len() && i < end; i += step {
		r, err := s.next.Get(value.Index(i).Interface())
		if err != nil {
			continue
		}
		result = append(result, r.data)
	}
	return result, nil
}
