package jtdinfer

type Inferrer struct {
	Inference *InferredSchema
	// TODO: Add hints
}

func NewInferrer() *Inferrer {
	return &Inferrer{
		Inference: NewInferredSchema(),
	}
}

func (i *Inferrer) Infer(value any) *Inferrer {
	return &Inferrer{
		Inference: i.Inference.Infer(value),
	}
}

func (i *Inferrer) IntoSchema() Schema {
	return i.Inference.IntoSchema()
}
