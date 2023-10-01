package main

import (
	"encoding/json"

	"github.com/bombsimon/jtdinfer"
)

func main() {
	schema := jtdinfer.
		NewInferrer(jtdinfer.WithoutHints()).
		Infer("my-string").
		IntoSchema()

	j, _ := json.MarshalIndent(schema, "", "  ")
	print(string(j))
}
