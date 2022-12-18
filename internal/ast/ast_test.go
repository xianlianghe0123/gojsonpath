package ast

type Book struct {
	Category string  `json:"category"`
	Author   string  `json:"author"`
	Title    string  `json:"title"`
	Isbn     string  `json:"isbn"`
	Price    float64 `json:"price"`
}

var data = map[string]interface{}{
	"store": map[string]interface{}{
		"book": []interface{}{
			&Book{
				Category: "reference",
				Author:   "Nigel Rees",
				Title:    "Sayings of the Century",
				Price:    8.95,
			},
			map[string]interface{}{
				"category": "fiction",
				"author":   "Evelyn Waugh",
				"title":    "Sword of Honour",
				"price":    12.99,
			},
			&Book{
				Category: "fiction",
				Author:   "Herman Melville",
				Title:    "Moby Dick",
				Isbn:     "0-553-21311-3",
				Price:    8.99,
			},
			map[string]interface{}{
				"category": "fiction",
				"author":   "J. R. R. Tolkien",
				"title":    "The Lord of the Rings",
				"isbn":     "0-395-19395-8",
				"price":    22.99,
			},
		},
		"bicycle": map[string]interface{}{
			"color": "red",
			"price": 19.95,
		},
	},
}
