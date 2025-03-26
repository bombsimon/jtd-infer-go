package main

import (
	"encoding/json"
	"fmt"

	jtdinfer "github.com/bombsimon/jtd-infer-go"
)

func main() {
	inferSimpleValue()
	inferMultipleStringRows()
	inferMap()
	inferManualUnmarshal()
	inferWithHints()
}

func inferSimpleValue() {
	schema := jtdinfer.
		NewInferrer(jtdinfer.WithoutHints()).
		Infer("my-string").
		IntoSchema()

	j, _ := json.MarshalIndent(schema, "", "  ")
	fmt.Println(string(j))
	fmt.Println()
}

func inferMap() {
	schema := jtdinfer.
		NewInferrer(jtdinfer.WithoutHints()).
		Infer(map[string]any{
			"age":  52,
			"name": "Joe",
		}).
		IntoSchema()

	j, _ := json.MarshalIndent(schema, "", "  ")
	fmt.Println(string(j))
	fmt.Println()
}

func inferMultipleStringRows() {
	rows := []string{
		`{"name":"Joe", "age": 52, "something_optional": true, "something_nullable": 1.1}`,
		`{"name":"Jane", "age": 48, "something_nullable": null}`,
	}
	schema := jtdinfer.
		InferStrings(rows, jtdinfer.WithoutHints()).
		IntoSchema()

	j, _ := json.MarshalIndent(schema, "", "  ")
	fmt.Println(string(j))
	fmt.Println()
}

func inferManualUnmarshal() {
	var m map[string]any
	json.Unmarshal([]byte(`{"name": "Jon", "age": 52}`), &m)

	schema := jtdinfer.
		NewInferrer(jtdinfer.WithoutHints()).
		Infer(m).
		IntoSchema()

	j, _ := json.MarshalIndent(schema, "", "  ")
	fmt.Println(string(j))
	fmt.Println()
}

func inferWithHints() {
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
	hints := jtdinfer.Hints{
		DefaultNumType: jtdinfer.NumTypeUint32,
		Enums:          jtdinfer.NewHintSet().Add([]string{"work", "department"}),
		Values:         jtdinfer.NewHintSet().Add([]string{"values"}),
		Discriminator:  jtdinfer.NewHintSet().Add([]string{"discriminator", "-", "type"}),
	}

	schema := jtdinfer.InferStrings(rows, hints).IntoSchema()
	j, _ := json.MarshalIndent(schema, "", "  ")
	fmt.Println(string(j))
	fmt.Println()
}
