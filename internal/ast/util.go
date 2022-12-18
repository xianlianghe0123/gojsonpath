package ast

import (
	"reflect"
	"strings"
)

func getFieldKey(sf reflect.StructField) (key string, omitempty bool) {
	if !sf.IsExported() {
		return "", false
	}
	key = sf.Name
	tag, ok := sf.Tag.Lookup("json")
	if !ok {
		return key, false
	}
	for i, s := range strings.Split(tag, ",") {
		if i == 0 {
			key = s
			continue
		}
		if s == "omitempty" {
			omitempty = true
			break
		}
	}
	return key, omitempty
}
