package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	jtdinfer "github.com/bombsimon/jtd-infer-go"
)

func main() {
	inferrer := jtdinfer.NewInferrer(jtdinfer.WithoutHints())
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() && scanner.Text() != "" {
		var toInfer any
		if err := json.Unmarshal(scanner.Bytes(), &toInfer); err != nil {
			fmt.Printf("failed to infer string: %s\n", err)
			os.Exit(1)
		}

		inferrer = inferrer.Infer(toInfer)
	}

	schema := inferrer.IntoSchema()
	j, _ := json.MarshalIndent(schema, "", "  ") //nolint:errchkjson
	fmt.Println(string(j))
}
