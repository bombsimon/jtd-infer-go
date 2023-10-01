package jtdinfer

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
func (i *Inferrer) IntoSchema() Schema {
	return i.Inference.IntoSchema()
}
