package jtdinfer

import (
	"encoding/json"
	"strconv"
)

// Inferrer represents the `InferredSchema` with its state combined with the
// hints used when inferring.
type Inferrer struct {
	Inference *InferredSchema
	Hints     *Hints
}

// NewInferrer will create a new inferrer with a default `InferredSchema`.
func NewInferrer(hints *Hints) *Inferrer {
	return &Inferrer{
		Inference: NewInferredSchema(),
		Hints:     hints,
	}
}

// Infer will infer the schema.
func (i *Inferrer) Infer(value any) *Inferrer {
	return &Inferrer{
		Inference: i.Inference.Infer(value, i.Hints),
		Hints:     i.Hints,
	}
}

// IntoSchema will convert the `InferredSchema` into a final `Schema`.
func (i *Inferrer) IntoSchema() Schema {
	return i.Inference.IntoSchema(i.Hints)
}

// InferStrings accepts a slice of strings and will try to JSON unmarshal each
// row to the type that the first row looks like. If an error occurs the
// inferrer will return with the state it had when the error occurred.
// If you already have the type of your data such as a slice of numbers or a map
// of strings you can pass them directly to `Infer`. This is just a convenience
// method if all you got is strings.
func InferStrings(rows []string, hints *Hints) *Inferrer {
	inferrer := NewInferrer(hints)
	if len(rows) == 0 {
		return inferrer
	}

	var (
		firstRow   = rows[0]
		getToInfer func() any
	)

	switch {
	case isBool(firstRow):
		getToInfer = func() any { return false }
	case isObject(firstRow):
		getToInfer = func() any { return make(map[string]any) }
	case isArray(firstRow):
		getToInfer = func() any { return make([]any, 0) }
	case isNumber(firstRow):
		getToInfer = func() any { return 0.0 }
	default:
		getToInfer = func() any { return "" }
	}

	for _, row := range rows {
		toInfer := getToInfer()
		if err := json.Unmarshal([]byte(row), &toInfer); err != nil {
			return inferrer
		}

		inferrer = inferrer.Infer(toInfer)
	}

	return inferrer
}

func isBool(value string) bool {
	return value == "true" || value == "false"
}

func isObject(value string) bool {
	var m map[string]any
	return json.Unmarshal([]byte(value), &m) == nil
}

func isArray(value string) bool {
	var a []any
	return json.Unmarshal([]byte(value), &a) == nil
}

func isNumber(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}
