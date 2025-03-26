package jtdinfer

import (
	"fmt"
	"math"
	"strconv"
	"testing"
	"time"

	jtd "github.com/jsontypedef/json-typedef-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInfer(t *testing.T) {
	for _, tc := range []struct {
		value          any
		expectedSchema Schema
	}{
		{
			value:          "52",
			expectedSchema: Schema{Type: jtd.TypeString},
		},
		{
			value:          52,
			expectedSchema: Schema{Type: jtd.TypeUint8},
		},
		{
			value:          nil,
			expectedSchema: Schema{Nullable: true},
		},
		{
			value: map[string]any{
				"name": "Joe",
				"age":  52,
			},
			expectedSchema: Schema{
				Properties: map[string]Schema{
					"name": {Type: jtd.TypeString},
					"age":  {Type: jtd.TypeUint8},
				},
			},
		},
		{
			value: []int{1, 3, 5},
			expectedSchema: Schema{
				Elements: &Schema{Type: jtd.TypeUint8},
			},
		},
	} {
		t.Run(fmt.Sprintf("infer_%T", tc.value), func(t *testing.T) {
			gotSchema := NewInferrer(WithoutHints()).Infer(tc.value).IntoSchema()
			assert.EqualValues(t, tc.expectedSchema, gotSchema)
		})
	}
}

func TestInferString(t *testing.T) {
	cases := []struct {
		description    string
		values         []string
		expectedSchema Schema
	}{
		{
			description: "boolean true value",
			values:      []string{"true"},
			expectedSchema: Schema{
				Type: jtd.TypeBoolean,
			},
		},
		{
			description: "boolean false value",
			values:      []string{"false"},
			expectedSchema: Schema{
				Type: jtd.TypeBoolean,
			},
		},
		{
			description: "object",
			values:      []string{`{"name":"Joe"}`},
			expectedSchema: Schema{
				Properties: map[string]Schema{
					"name": {
						Type: jtd.TypeString,
					},
				},
			},
		},
		{
			description: "object first is null",
			values:      []string{`{"name":null}`, `{"name":"Joe"}`},
			expectedSchema: Schema{
				Properties: map[string]Schema{
					"name": {
						Type:     jtd.TypeString,
						Nullable: true,
					},
				},
			},
		},
		{
			description: "object with array",
			values: []string{
				`{"name": "Joe", "age": 42, "hobbies": ["code", "animals"]}`,
			},
			expectedSchema: Schema{
				Properties: map[string]Schema{
					"name":    {Type: jtd.TypeString},
					"age":     {Type: jtd.TypeUint8},
					"hobbies": {Elements: &Schema{Type: jtd.TypeString}},
				},
			},
		},
		{
			description: "array",
			values:      []string{`[1, 2, 3]`},
			expectedSchema: Schema{
				Elements: &Schema{
					Type: jtd.TypeUint8,
				},
			},
		},
		{
			description: "unsigned integer",
			values:      []string{"1"},
			expectedSchema: Schema{
				Type: jtd.TypeUint8,
			},
		},
		{
			description: "signed integer",
			values:      []string{"-1"},
			expectedSchema: Schema{
				Type: jtd.TypeInt8,
			},
		},
		{
			description: "signed max integer",
			values:      []string{strconv.Itoa(math.MinInt32)},
			expectedSchema: Schema{
				Type: jtd.TypeInt32,
			},
		},
		{
			description: "float without fraction",
			values:      []string{"1.0"},
			expectedSchema: Schema{
				Type: jtd.TypeUint8,
			},
		},
		{
			description: "positive float",
			values:      []string{"1.1"},
			expectedSchema: Schema{
				Type: jtd.TypeFloat64,
			},
		},
		{
			description: "negative float",
			values:      []string{"-1.1"},
			expectedSchema: Schema{
				Type: jtd.TypeFloat64,
			},
		},
		{
			description: "string",
			values:      []string{`"string"`},
			expectedSchema: Schema{
				Type: jtd.TypeString,
			},
		},
		{
			description: "number in string is still string",
			values:      []string{`"2.2"`},
			expectedSchema: Schema{
				Type: jtd.TypeString,
			},
		},
		{
			description: "timestamp",
			values:      []string{fmt.Sprintf(`"%s"`, time.Now().Format(time.RFC3339))},
			expectedSchema: Schema{
				Type: jtd.TypeTimestamp,
			},
		},
		{
			description: "null",
			values:      []string{"null"},
			expectedSchema: Schema{
				Nullable: true,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			gotSchema := InferStrings(tc.values, WithoutHints()).IntoSchema()
			assert.EqualValues(t, tc.expectedSchema, gotSchema)
		})
	}
}

func TestInferrerWithEnumHints(t *testing.T) {
	hints := Hints{
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

func TestInferWithValuesHints(t *testing.T) {
	hints := Hints{
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

func TestInferWithDiscriminatorHints(t *testing.T) {
	hints := Hints{
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

func BenchmarkInferOneRowNoMissingHints(b *testing.B) {
	rows := generateRows(1)
	emptyHints := WithoutHints()

	for n := 0; n < b.N; n++ {
		InferStrings(rows, emptyHints)
	}
}

func BenchmarkInferThousandRowsNoMissingHints(b *testing.B) {
	rows := generateRows(1000)
	emptyHints := WithoutHints()

	for n := 0; n < b.N; n++ {
		InferStrings(rows, emptyHints)
	}
}

func BenchmarkInferSimpleString(b *testing.B) {
	inferrer := NewInferrer(WithoutHints())

	for n := 0; n < b.N; n++ {
		inferrer.Infer("string")
	}
}

func generateRows(n int) []string {
	row := `{"name":"bench", "speed":100.2}`
	rows := []string{}

	for i := 0; i < n; i++ {
		rows = append(rows, row)
	}

	return rows
}
