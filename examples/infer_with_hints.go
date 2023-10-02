package main

import (
	"encoding/json"

	jtdinfer "github.com/bombsimon/jtd-infer-go"
)

func main() {
	rows := []string{
		`{
			"name":"Joe",
			"age":52,
			"work":{"department": "sales"},
			"values":{"x": [1, 2, 3], "y": [4, 5, 6], "z": [7, 8, 9]},
			"discriminator":[{"type":"s", "value":"foo"},{"type":"n", "value":3.14}]
		}`,
		`{
			"name":"Jane",
			"age":48,
			"work":{"department": "engineering"},
			"values":{"x": [1, 2, 3], "y": [4, 5, 6], "z": [7, 8, -2000]},
			"discriminator":[{"type":"s", "value":"foo"},{"type":"n", "value":3.14}]
		}`,
	}
	hints := &jtdinfer.Hints{
		DefaultNumType: jtdinfer.NumTypeUint32,
		Enums:          jtdinfer.NewHintSet().Add([]string{"work", "department"}),
		Values:         jtdinfer.NewHintSet().Add([]string{"values"}),
		Discriminator:  jtdinfer.NewHintSet().Add([]string{"discriminator", "-", "type"}),
	}

	schema := jtdinfer.InferStrings(rows, hints).IntoSchema()
	j, _ := json.MarshalIndent(schema, "", "  ")
	print(string(j))
}
