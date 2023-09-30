package jtdinfer

import (
	"math"

	jtd "github.com/jsontypedef/json-typedef-go"
)

type InferredNumber struct {
	Min       float64
	Max       float64
	IsInteger bool
}

func NewNumber() *InferredNumber {
	return &InferredNumber{
		IsInteger: true,
	}
}

func (i *InferredNumber) Infer(n float64) *InferredNumber {
	return &InferredNumber{
		Min:       math.Min(i.Min, n),
		Max:       math.Max(i.Max, n),
		IsInteger: i.IsInteger && float64(int(n)) == n,
	}
}

func (i *InferredNumber) IntoType(defaultType minMax) jtd.Type {
	if i.ContainedBy(defaultType) {
		return defaultType.typ
	}

	mm := []minMax{
		{typ: jtd.TypeUint8, min: 0, max: math.MaxUint8},
		{typ: jtd.TypeInt8, min: math.MinInt8, max: math.MaxInt8},
		{typ: jtd.TypeUint16, min: 0, max: math.MaxUint16},
		{typ: jtd.TypeInt16, min: math.MinInt16, max: math.MaxInt16},
		{typ: jtd.TypeUint32, min: 0, max: math.MaxUint32},
		{typ: jtd.TypeInt32, min: math.MinInt32, max: math.MaxInt32},
	}

	for _, v := range mm {
		if i.ContainedBy(v) {
			return v.typ
		}
	}

	return jtd.TypeFloat64
}

type minMax struct {
	typ jtd.Type
	min float64
	max float64
}

func (i *InferredNumber) ContainedBy(v minMax) bool {
	if !i.IsInteger && v.typ != jtd.TypeFloat32 && v.typ != jtd.TypeFloat64 {
		return false
	}

	return v.min <= i.Min && v.max >= i.Max
}
