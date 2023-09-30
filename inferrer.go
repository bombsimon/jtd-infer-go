package jtdinfer

// Inferrer represents the `InferredSchema` with its state combined with the
// hints used when inferring.
type Inferrer struct {
	Inference *InferredSchema
	// TODO: Add hints
}

// NewInferred will create a new inferrer with a default `InferredSchema`.
func NewInferrer() *Inferrer {
	return &Inferrer{
		Inference: NewInferredSchema(),
	}
}

// Infer will infer the schema.
func (i *Inferrer) Infer(value any) *Inferrer {
	return &Inferrer{
		Inference: i.Inference.Infer(value),
	}
}

// IntoSchema will convert the `InferredSchema` into a final `Schema`.
func (i *Inferrer) IntoSchema() Schema {
	return i.Inference.IntoSchema()
}
