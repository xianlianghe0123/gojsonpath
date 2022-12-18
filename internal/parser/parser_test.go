package parser

import (
	"testing"
)

func TestParser_Parse(t *testing.T) {
	cases := []struct {
		jsonPath    string
		expectation string
	}{
		{`$.a['b',"c"]`, `$["a"]["b","c"]`},
		{`$.a[1:5:3,7,8]`, `$["a"][1:5:3,7,8]`},
		{`$.a.b.c`, `$["a"]["b"]["c"]`},
		{`$. $a`, `$[" $a"]`},
		{`$.['a\'a', "b\"b"]`, `$["a'a","b\"b"]`},

		{`$....a`, ``},
		{`$[1`, ``},
	}
	for _, c := range cases {
		ast, err := NewParser(c.jsonPath).Parse()
		if err != nil {
			t.Logf("Case %s err: %+v", c.jsonPath, err)
		}
		cur := ast.String()
		if cur != c.expectation {
			t.Errorf("Case %s expected:%s, current:%s\n", c.jsonPath, c.expectation, cur)
		}
	}
}
