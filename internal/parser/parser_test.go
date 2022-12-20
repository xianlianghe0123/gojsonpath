package parser

import (
	"testing"
)

func TestParser_Parse(t *testing.T) {
	cases := []struct {
		jsonPath    string
		hasErr      bool
		expectation string
	}{
		{`$.a['b',"c"]`, false, `$["a"]["b","c"]`},
		{`$.[:]`, false, `$[::]`},
		{`$.[1:]`, false, `$[1::]`},
		{`$.[:3:]`, false, `$[:3:]`},
		{`$.[::2]`, false, `$[::2]`},
		{`$.[:,2]`, false, `$[::,2]`},
		{`$.[:,2,:]`, false, `$[::,2,::]`},
		{`$.a[1:5:3,7,8].a`, false, `$["a"][1:5:3,7,8]`},
		{`$.a.b.c`, false, `$["a"]["b"]["c"]`},
		{`$. $a`, false, `$[" $a"]`},
		{`$.['a\'a', "b\"b"]`, false, `$["a'a","b\"b"]`},

		{`$....a`, true, ``},
		{`$[1`, true, ``},
	}
	for _, c := range cases {
		ast, err := NewParser(c.jsonPath).Parse()
		if err != nil {
			t.Logf("Case %s err: %+v", c.jsonPath, err)
			if !c.hasErr {
				t.Failed()
			}
			continue
		}
		if c.hasErr {
			t.Failed()
		}
		cur := ast.String()
		if cur != c.expectation {
			t.Errorf("Case %s expected:%s, current:%s\n", c.jsonPath, c.expectation, cur)
		}
	}
}
