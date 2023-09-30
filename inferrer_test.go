package jtdinfer

import (
	"encoding/json"
	"testing"

	jtd "github.com/jsontypedef/json-typedef-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJTDSimpleInfer(t *testing.T) {
	rows := []string{
		`{ "name": "Joe", "age": 42 }`,
	}

	inferrer := NewInferrer()

	for _, row := range rows {
		rowAsJSON := make(map[string]any, 0)
		require.NoError(t, json.Unmarshal([]byte(row), &rowAsJSON))
		inferrer = inferrer.Infer(rowAsJSON)
	}

	expectedSchema := Schema{
		Properties: map[string]Schema{
			"name": {Type: jtd.TypeString},
			"age":  {Type: jtd.TypeUint8},
		},
	}
	gotSchema := inferrer.IntoSchema()

	assert.EqualValues(t, expectedSchema, gotSchema)
}
