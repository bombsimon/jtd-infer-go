package jtdinfer

import (
	"testing"

	jtd "github.com/jsontypedef/json-typedef-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJTDInfer(t *testing.T) {
	rows := []string{
		`{"name": "Joe", "age": 42, "hobbies": ["code", "animals"]}`,
	}

	expectedSchema := Schema{
		Properties: map[string]Schema{
			"name":    {Type: jtd.TypeString},
			"age":     {Type: jtd.TypeUint8},
			"hobbies": {Elements: &Schema{Type: jtd.TypeString}},
		},
	}
	gotSchema := InferStrings(rows, NewHints()).IntoSchema()

	assert.EqualValues(t, expectedSchema, gotSchema)
}

func TestJTDInferrerWithEnumHints(t *testing.T) {
	hints := &Hints{
		Enums: NewHintSet().
			Add([]string{"name"}).
			Add([]string{"address", "city"}),
	}

	rows := []string{
		`{"address": {"city": "Stockholm"}, "name": "Joe", "age": 42}`,
		`{"address": {"city": "Umeå"}, "name": "Labero", "age": 42}`,
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
	gotSchema := InferStrings(rows, hints).IntoSchema()

	// We check that we got the same elements in our enum first and then we
	// delete it since the order is unreliable due to being a map.
	require.ElementsMatch(
		t,
		expectedSchema.Properties["name"].Enum,
		gotSchema.Properties["name"].Enum,
	)

	delete(expectedSchema.Properties, "name")
	delete(gotSchema.Properties, "name")

	require.ElementsMatch(
		t,
		expectedSchema.Properties["address"].Properties["city"].Enum,
		gotSchema.Properties["address"].Properties["city"].Enum,
	)

	delete(expectedSchema.Properties, "address")
	delete(gotSchema.Properties, "address")

	assert.EqualValues(t, expectedSchema, gotSchema)
}

func TestJTDInferWithValuesHints(t *testing.T) {
	hints := &Hints{
		Values: NewHintSet().Add([]string{}),
	}

	rows := []string{
		`{"x": [1, 2, 3], "y": [4, 5, 6], "z": [7, 8, 9]}`,
		`{"x": [1, 2, 3], "y": [4, 5, -600], "z": [7, 8, 9]}`,
	}

	expectedSchema := Schema{
		Values: &Schema{
			Elements: &Schema{
				Type: jtd.TypeInt16,
			},
		},
	}
	gotSchema := InferStrings(rows, hints).IntoSchema()

	assert.EqualValues(t, expectedSchema, gotSchema)
}

func TestJTDInferWithDiscriminatorHints(t *testing.T) {
	hints := &Hints{
		Discriminator: NewHintSet().Add([]string{"-", "type"}),
	}

	rows := []string{
		`[{"type": "s", "value": "foo"},{"type": "n", "value": 3.14}]`,
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
	gotSchema := InferStrings(rows, hints).IntoSchema()

	assert.EqualValues(t, expectedSchema, gotSchema)
}
