package main

import (
	"encoding/json"

	"github.com/bombsimon/jtdinfer"
)

func main() {
	rows := []string{
		`{"name":"Joe", "age": 52, "something_optional": true, "something_nullable": 1.1}`,
		`{"name":"Jane", "age": 48, "something_nullable": null}`,
	}
	schema := jtdinfer.
		InferStrings(rows, jtdinfer.WithoutHints()).
		IntoSchema(jtdinfer.WithoutHints())

	j, _ := json.MarshalIndent(schema, "", "  ")
	print(string(j))
}
