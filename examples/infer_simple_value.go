package main

import (
	"encoding/json"

	"github.com/bombsimon/jtdinfer"
)

func main() {
	schema := jtdinfer.
		NewInferrer(jtdinfer.WithoutHints()).
		Infer("my-string").
		IntoSchema(jtdinfer.WithoutHints())

	j, _ := json.MarshalIndent(schema, "", "  ")
	print(string(j))
}
