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

// NewInferred will create a new inferrer with a default `InferredSchema`.
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
func (i *Inferrer) IntoSchema(hints Hints) Schema {
	return i.Inference.IntoSchema(hints)
}

// InferStrings accepts a slice of strings and will convert them to either a
// `map[string]any` or []any` and run inference on all the rows. If any of the
// rows are not valid JSON object or list, the inference up to that point is
// returned.
//
// If you need to infer simple values like strings or integers they can be
// passed directly to `Infer`.
func InferStrings(rows []string, hints Hints) *Inferrer {
	inferrer := NewInferrer(hints)

	for _, row := range rows {
		var toInfer any = make(map[string]any, 0)
		if err := json.Unmarshal([]byte(row), &toInfer); err != nil {
			toInfer = make([]any, 0)
			if err := json.Unmarshal([]byte(row), &toInfer); err != nil {
				return inferrer
			}
		}

		inferrer = inferrer.Infer(toInfer)
	}

	return inferrer
}
