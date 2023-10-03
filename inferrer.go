package jtdinfer

import (
	"encoding/json"
)

// Inferrer represents the `InferredSchema` with its state combined with the
// hints used when inferring.
type Inferrer struct {
	Inference *InferredSchema
	Hints     Hints
}

// NewInferrer will create a new inferrer with a default `InferredSchema`.
func NewInferrer(hints Hints) *Inferrer {
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
// row. If an error occurs the inferrer will return with the state it had when
// the error occurred. If you already have the type of your data such as a slice
// of numbers or a map of strings you can pass them directly to `Infer`. This is
// just a convenience method if all you got is strings.
func InferStrings(rows []string, hints Hints) *Inferrer {
	inferrer := NewInferrer(hints)
	if len(rows) == 0 {
		return inferrer
	}

	for _, row := range rows {
		var toInfer any
		if err := json.Unmarshal([]byte(row), &toInfer); err != nil {
			return inferrer
		}

		inferrer = inferrer.Infer(toInfer)
	}

	return inferrer
}
