# jsonpath
Introduction: https://goessner.net/articles/JsonPath/


## Syntax Support
|         JSONPath         | Description            | Support |
|:------------------------:|:-----------------------|:-------:|
|           `$`            | return root            |    ✅    |
|           `.`            | child                  |    ✅    |
|           `..`           | recursive              |    ✅    |
|       `*` or `[*]`       | all                    |    ✅    |
|          `[0]`           | index                  |    ✅    |
|          `[-1]`          | negative index         |    ✅    |
|    `[start:end:step]`    | slice                  |    ✅    |
| `[start:end:step,0,...]` | union: index and slice |    ✅    |
| `['field'] or ["field"]` | field                  |    ✅    |
|     `['field1',...]`     | union: fields          |    ✅    |
|         `[?()]`          | filter                 |    ❌    |
|          `[()]`          | script expression      |    ❌    |

`*` or `..` order:
- `map`：random order (depend `reflect.MapRange`)
- `struct`：order by struct fields defined order

## Example
```json
{
  "store": {
    "book": [
      {
        "category": "reference",
        "author": "Nigel Rees",
        "title": "Sayings of the Century",
        "price": 8.95
      },
      {
        "category": "fiction",
        "author": "Evelyn Waugh",
        "title": "Sword of Honour",
        "price": 12.99
      },
      {
        "category": "fiction",
        "author": "Herman Melville",
        "title": "Moby Dick",
        "isbn": "0-553-21311-3",
        "price": 8.99
      },
      {
        "category": "fiction",
        "author": "J. R. R. Tolkien",
        "title": "The Lord of the Rings",
        "isbn": "0-395-19395-8",
        "price": 22.99
      }
    ],
    "bicycle": {
      "color": "red",
      "price": 19.95
    }
  }
}
```
|               JSONPath               | Result                                                                                                                                                                                                                               |
|:------------------------------------:|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|                 `$`                  | root<br/>`{"store":{...}}`                                                                                                                                                                                                           |
|       `$.store.book[*].author`       | the authors of all books in the store<br/>`["Nigel Rees","Evelyn Waugh","Herman Melville","J. R. R. Tolkien"]`                                                                                                                       |
| `$.store.book[*].['author',"price"]` | the authors and price of all books in the store<br/>`["Nigel Rees",8.95,"Evelyn Waugh",12.99,"Herman Melville",8.99,"J. R. R. Tolkien",22.99]`                                                                                       |
|             `$..author`              | all authors<br/>`["Nigel Rees","Evelyn Waugh","Herman Melville","J. R. R. Tolkien"]`                                                                                                                                                 |
|             `$.store.*`              | all things in store, which are some books and a red bicycle<br/>`["book":[...],"bicycle":{"color":"red","price":19.95}]`                                                                                                             |
|          	`$.store..price`           | the price of everything in the store<br/>`[8.95,12.99,8.99,22.99,19.95]`                                                                                                                                                             |
|             `$..book[2]`             | the third book<br/>`{"category":"fiction","author":"Herman Melville","title":"Moby Dick","isbn":"0-553-21311-3","price":8.99}`                                                                                                       |
|            `$..book[-1]`             | the last book in order<br/>`{"category":"fiction","author":"J. R. R. Tolkien","title":"The Lord of the Rings","isbn":"0-395-19395-8","price":22.99}`                                                                                 |
|            `$..book[:2]`             | the first two books<br/>`[{"category":"reference","author":"Nigel Rees","title":"Sayings of the Century","price":8.95},{"category":"fiction","author":"Evelyn Waugh","title":"Sword of Honour","price":12.99}]`                      |
|           `$..book[:2,3]`            | the first two books and the fourth book<br/>`[{"category":"reference","author":"Nigel Rees","title":"Sayings of the Century","price":8.95},{"category":"fiction","author":"Evelyn Waugh","title":"Sword of Honour","price":12.99},{"category":"fiction","author":"J. R. R. Tolkien","title":"The Lord of the Rings","isbn":"0-395-19395-8","price":22.99}]` |