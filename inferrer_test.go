package jtdinfer

import (
	"encoding/json"
	"testing"

	jtd "github.com/jsontypedef/json-typedef-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJTDInfer(t *testing.T) {
	rows := []string{
		`{"name": "Joe", "age": 42, "hobbies": ["code", "animals"]}`,
	}

	inferrer := NewInferrer(Hints{})

	for _, row := range rows {
		rowAsJSON := make(map[string]any, 0)
		require.NoError(t, json.Unmarshal([]byte(row), &rowAsJSON))
		inferrer = inferrer.Infer(rowAsJSON)
	}

	expectedSchema := Schema{
		Properties: map[string]Schema{
			"name":    {Type: jtd.TypeString},
			"age":     {Type: jtd.TypeUint8},
			"hobbies": {Elements: &Schema{Type: jtd.TypeString}},
		},
	}
	gotSchema := inferrer.IntoSchema()

	assert.EqualValues(t, expectedSchema, gotSchema)
}

func TestJTDInferrerWithHints(t *testing.T) {
	hints := Hints{
		Enums: HintSet{
			Values: [][]string{
				{"name"},
				{"address", "city"},
			},
		},
	}

	rows := []string{
		`{"address": {"city": "Stockholm"}, "name": "Joe", "age": 42}`,
		`{"address": {"city": "Umeå"}, "name": "Labero", "age": 42}`,
	}

	inferrer := NewInferrer(hints)

	for _, row := range rows {
		rowAsJSON := make(map[string]any, 0)
		require.NoError(t, json.Unmarshal([]byte(row), &rowAsJSON))
		inferrer = inferrer.Infer(rowAsJSON)
	}

	expectedSchema := Schema{
		Properties: map[string]Schema{
			"name": {Enum: []string{"Joe", "Labero"}},
			"age":  {Type: jtd.TypeUint8},
			"address": {
				Properties: map[string]Schema{
					"city": {Enum: []string{"Stockholm", "Umeå"}},
				},
			},
		},
	}
	gotSchema := inferrer.IntoSchema()

	assert.EqualValues(t, expectedSchema, gotSchema)
}

func TestJTDInferWithDiscriminatorHints(t *testing.T) {
	hints := Hints{
		Discriminator: HintSet{
			Values: [][]string{
				{"-", "type"},
			},
		},
	}

	rows := []string{
		`[{"type": "s", "value": "foo"},{"type": "n", "value": 3.14}]`,
	}

	inferrer := NewInferrer(hints)

	for _, row := range rows {
		rowAsJSON := make([]any, 0)
		require.NoError(t, json.Unmarshal([]byte(row), &rowAsJSON))
		inferrer = inferrer.Infer(rowAsJSON)
	}

	expectedSchema := Schema{
		Elements: &Schema{
			Discriminator: "type",
			Mapping: map[string]Schema{
				"s": {
					Properties: map[string]Schema{
						"value": {Type: jtd.TypeString},
					},
				},
				"n": {
					Properties: map[string]Schema{
						"value": {Type: jtd.TypeFloat64},
					},
				},
			},
		},
	}
	gotSchema := inferrer.IntoSchema()

	assert.EqualValues(t, expectedSchema, gotSchema)
}
