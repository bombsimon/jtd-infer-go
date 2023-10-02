package main

import (
	"encoding/json"

	jtdinfer "github.com/bombsimon/jtd-infer-go"
)

func main() {
	schema := jtdinfer.
		NewInferrer(jtdinfer.NewHints()).
		Infer("my-string").
		IntoSchema()

	j, _ := json.MarshalIndent(schema, "", "  ")
	print(string(j))
}
