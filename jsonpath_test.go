package jsonpath

import (
	"encoding/json"
	"reflect"
	"testing"
)

type Book struct {
	Category string  `json:"category"`
	Author   string  `json:"author"`
	Title    string  `json:"title"`
	Isbn     string  `json:"isbn,omitempty"`
	Price    float64 `json:"price"`
}

type Bicycle struct {
	Color string  `json:"color"`
	Price float64 `json:"price"`
}

var data = map[string]map[string]interface{}{
	"store": {
		"book": []*Book{
			{
				Category: "reference",
				Author:   "Nigel Rees",
				Title:    "Sayings of the Century",
				Price:    8.95,
			},
			{
				Category: "fiction",
				Author:   "Evelyn Waugh",
				Title:    "Sword of Honour",
				Price:    12.99,
			},
			{
				Category: "fiction",
				Author:   "Herman Melville",
				Title:    "Moby Dick",
				Price:    8.99,
				Isbn:     "0-553-21311-3",
			},
			{
				Category: "fiction",
				Author:   "J. R. R. Tolkien",
				Title:    "The Lord of the Rings",
				Isbn:     "0-395-19395-8",
				Price:    22.99,
			},
		},
		"bicycle": &Bicycle{
			Color: "red",
			Price: 19.95,
		},
	},
}

func TestGet(t *testing.T) {
	cases := []struct {
		jsonPath    string
		expectation string
	}{
		{`$`, `{"store":{"book":[{"category":"reference","author":"Nigel Rees","title":"Sayings of the Century","price":8.95},{"category":"fiction","author":"Evelyn Waugh","title":"Sword of Honour","price":12.99},{"category":"fiction","author":"Herman Melville","title":"Moby Dick","isbn":"0-553-21311-3","price":8.99},{"category":"fiction","author":"J. R. R. Tolkien","title":"The Lord of the Rings","isbn":"0-395-19395-8","price":22.99}],"bicycle":{"color":"red","price":19.95}}}`},
		{`$.*`, `[{"book":[{"category":"reference","author":"Nigel Rees","title":"Sayings of the Century","price":8.95},{"category":"fiction","author":"Evelyn Waugh","title":"Sword of Honour","price":12.99},{"category":"fiction","author":"Herman Melville","title":"Moby Dick","isbn":"0-553-21311-3","price":8.99},{"category":"fiction","author":"J. R. R. Tolkien","title":"The Lord of the Rings","isbn":"0-395-19395-8","price":22.99}],"bicycle":{"color":"red","price":19.95}}]`},
		{`$.store.book[*].author`, `["Nigel Rees","Evelyn Waugh","Herman Melville","J. R. R. Tolkien"]`},
		{`$.store.book[*].['author',"price"]`, `["Nigel Rees",8.95,"Evelyn Waugh",12.99,"Herman Melville",8.99,"J. R. R. Tolkien",22.99]`},
		{`$..author`, `["Nigel Rees","Evelyn Waugh","Herman Melville","J. R. R. Tolkien"]`},
		{`$.store.*`, `[[{"category":"reference","author":"Nigel Rees","title":"Sayings of the Century","price":8.95},{"category":"fiction","author":"Evelyn Waugh","title":"Sword of Honour","price":12.99},{"category":"fiction","author":"Herman Melville","title":"Moby Dick","isbn":"0-553-21311-3","price":8.99},{"category":"fiction","author":"J. R. R. Tolkien","title":"The Lord of the Rings","isbn":"0-395-19395-8","price":22.99}],{"color":"red","price":19.95}]`},
		{"$.store..price", `[8.95,12.99,8.99,22.99,19.95]`},
		{`$..book[2]`, `[{"category":"fiction","author":"Herman Melville","title":"Moby Dick","isbn":"0-553-21311-3","price":8.99}]`},
		{`$..book[-1:]`, `[{"category":"fiction","author":"J. R. R. Tolkien","title":"The Lord of the Rings","isbn":"0-395-19395-8","price":22.99}]`},
		{`$..book[0,1]`, `[{"category":"reference","author":"Nigel Rees","title":"Sayings of the Century","price":8.95},{"category":"fiction","author":"Evelyn Waugh","title":"Sword of Honour","price":12.99}]`},
		{`$..book[:2]`, `[{"category":"reference","author":"Nigel Rees","title":"Sayings of the Century","price":8.95},{"category":"fiction","author":"Evelyn Waugh","title":"Sword of Honour","price":12.99}]`},
		{`$..book[:2,3]`, `[{"category":"reference","author":"Nigel Rees","title":"Sayings of the Century","price":8.95},{"category":"fiction","author":"Evelyn Waugh","title":"Sword of Honour","price":12.99},{"category":"fiction","author":"J. R. R. Tolkien","title":"The Lord of the Rings","isbn":"0-395-19395-8","price":22.99}]`},
		{`$..*`, `[{"book":[{"category":"reference","author":"Nigel Rees","title":"Sayings of the Century","price":8.95},{"category":"fiction","author":"Evelyn Waugh","title":"Sword of Honour","price":12.99},{"category":"fiction","author":"Herman Melville","title":"Moby Dick","isbn":"0-553-21311-3","price":8.99},{"category":"fiction","author":"J. R. R. Tolkien","title":"The Lord of the Rings","isbn":"0-395-19395-8","price":22.99}],"bicycle":{"color":"red","price":19.95}},[{"category":"reference","author":"Nigel Rees","title":"Sayings of the Century","price":8.95},{"category":"fiction","author":"Evelyn Waugh","title":"Sword of Honour","price":12.99},{"category":"fiction","author":"Herman Melville","title":"Moby Dick","isbn":"0-553-21311-3","price":8.99},{"category":"fiction","author":"J. R. R. Tolkien","title":"The Lord of the Rings","isbn":"0-395-19395-8","price":22.99}],{"color":"red","price":19.95},{"category":"reference","author":"Nigel Rees","title":"Sayings of the Century","price":8.95},{"category":"fiction","author":"Evelyn Waugh","title":"Sword of Honour","price":12.99},{"category":"fiction","author":"Herman Melville","title":"Moby Dick","isbn":"0-553-21311-3","price":8.99},{"category":"fiction","author":"J. R. R. Tolkien","title":"The Lord of the Rings","isbn":"0-395-19395-8","price":22.99},"reference","Nigel Rees","Sayings of the Century",8.95,"fiction","Evelyn Waugh","Sword of Honour",12.99,"fiction","Herman Melville","Moby Dick","0-553-21311-3",8.99,"fiction","J. R. R. Tolkien","The Lord of the Rings","0-395-19395-8",22.99,"red",19.95]`},
	}
	for _, c := range cases {
		d, err := Get(c.jsonPath, data)
		if err != nil {
			t.Logf("Case %q err: %+v", c.jsonPath, err)
		}
		var cur interface{}
		b, _ := json.Marshal(d)
		json.Unmarshal(b, &cur)
		var e interface{}
		json.Unmarshal([]byte(c.expectation), &e)

		if !reflect.DeepEqual(cur, e) {
			t.Errorf("Case %q, current:%s, expectation:%s\n", c.jsonPath, string(b), c.expectation)
		}

	}
}
