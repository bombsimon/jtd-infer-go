package jtdinfer

import (
	"math"
	"testing"

	jtd "github.com/jsontypedef/json-typedef-go"
	"github.com/stretchr/testify/assert"
)

func TestInferredNumberDefault(t *testing.T) {
	nn1 := NewNumber()
	nn2 := NewNumber().Infer(0.0).Infer(float64(math.MaxUint8))
	nn3 := NewNumber().Infer(float64(math.MinInt8)).Infer(float64(math.MaxInt8))
	nn4 := NewNumber().Infer(0.5)

	cases := []struct {
		inferredNumber *InferredNumber
		numType        NumType
		jtdType        jtd.Type
	}{
		// Defaults are honored.
		{nn1, NumTypeUint8, jtd.TypeUint8},
		{nn1, NumTypeInt8, jtd.TypeInt8},
		{nn1, NumTypeUint16, jtd.TypeUint16},
		{nn1, NumTypeInt16, jtd.TypeInt16},
		{nn1, NumTypeUint32, jtd.TypeUint32},
		{nn1, NumTypeInt32, jtd.TypeInt32},
		{nn1, NumTypeFloat32, jtd.TypeFloat32},
		{nn1, NumTypeFloat64, jtd.TypeFloat64},

		// Expand to limits of uint8.
		{nn2, NumTypeUint8, jtd.TypeUint8},
		{nn2, NumTypeInt8, jtd.TypeUint8},
		{nn2, NumTypeUint16, jtd.TypeUint16},
		{nn2, NumTypeInt16, jtd.TypeInt16},
		{nn2, NumTypeUint32, jtd.TypeUint32},
		{nn2, NumTypeInt32, jtd.TypeInt32},
		{nn2, NumTypeFloat32, jtd.TypeFloat32},
		{nn2, NumTypeFloat64, jtd.TypeFloat64},

		// Expand to limits of int8.
		{nn3, NumTypeUint8, jtd.TypeInt8},
		{nn3, NumTypeInt8, jtd.TypeInt8},
		{nn3, NumTypeUint16, jtd.TypeInt8},
		{nn3, NumTypeInt16, jtd.TypeInt16},
		{nn3, NumTypeUint32, jtd.TypeInt8},
		{nn3, NumTypeInt32, jtd.TypeInt32},
		{nn3, NumTypeFloat32, jtd.TypeFloat32},
		{nn3, NumTypeFloat64, jtd.TypeFloat64},

		// Test including a non-integer.
		{nn4, NumTypeUint8, jtd.TypeFloat64},
		{nn4, NumTypeInt8, jtd.TypeFloat64},
		{nn4, NumTypeUint16, jtd.TypeFloat64},
		{nn4, NumTypeInt16, jtd.TypeFloat64},
		{nn4, NumTypeUint32, jtd.TypeFloat64},
		{nn4, NumTypeInt32, jtd.TypeFloat64},
		{nn4, NumTypeFloat32, jtd.TypeFloat32},
		{nn4, NumTypeFloat64, jtd.TypeFloat64},
	}

	for _, tc := range cases {
		t.Run(string(tc.jtdType), func(t *testing.T) {
			assert.Equal(t, tc.jtdType, tc.inferredNumber.IntoType(tc.numType))
		})
	}
}
