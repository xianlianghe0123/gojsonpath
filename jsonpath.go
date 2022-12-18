package jsonpath

import (
	"encoding/json"
	"github.com/xianlianghe0123/jsonpath/internal/ast"
	"github.com/xianlianghe0123/jsonpath/internal/parser"
	"unsafe"
)

type Compiled struct {
	a ast.AST
}

func Compile(jsonPath string) (*Compiled, error) {
	a, err := parser.NewParser(jsonPath).Parse()
	if err != nil {
		return nil, err
	}
	return &Compiled{
		a: a,
	}, nil
}

func MustCompile(jsonPath string) *Compiled {
	c, err := Compile(jsonPath)
	if err != nil {
		panic(err)
	}
	return c
}

func (c *Compiled) Get(data interface{}) (interface{}, error) {
	return c.a.Get(data)
}

func (c *Compiled) GetBytes(dataBytes []byte) (interface{}, error) {
	var data interface{}
	err := json.Unmarshal(dataBytes, &data)
	if err != nil {
		return nil, err
	}
	return c.Get(data)
}

func (c *Compiled) GetString(dataStr string) (interface{}, error) {
	return c.GetBytes(*(*[]byte)(unsafe.Pointer(&dataStr)))
}

func Get(jsonPath string, data interface{}) (interface{}, error) {
	c, err := Compile(jsonPath)
	if err != nil {
		return nil, err
	}
	return c.Get(data)
}

func GetBytes(jsonPath string, dataBytes []byte) (interface{}, error) {
	c, err := Compile(jsonPath)
	if err != nil {
		return nil, err
	}
	return c.GetBytes(dataBytes)
}

func GetString(jsonPath string, dataStr string) (interface{}, error) {
	c, err := Compile(jsonPath)
	if err != nil {
		return nil, err
	}
	return c.GetString(dataStr)
}
