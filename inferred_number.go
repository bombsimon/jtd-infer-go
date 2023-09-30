package jtdinfer

import (
	"math"

	jtd "github.com/jsontypedef/json-typedef-go"
)

// NumType represents the type of number a number should be represented in the
// JTD.
type NumType uint8

const (
	NumTypeUint8 NumType = iota
	NumTypeInt8
	NumTypeUint16
	NumTypeInt16
	NumTypeUint32
	NumTypeInt32
	NumTypeFloat32
	NumTypeFloat64
)

// IsFloat returns true if the `NumType` is a float.
func (n NumType) IsFloat() bool {
	return n == NumTypeFloat32 || n == NumTypeFloat64
}

// AsRange returns the maximum and minimum value for a `NumType`.
func (n NumType) AsRange() (float64, float64) {
	switch n {
	case NumTypeUint8:
		return 0, math.MaxUint8
	case NumTypeInt8:
		return math.MinInt8, math.MaxInt8
	case NumTypeUint16:
		return 0, math.MaxUint16
	case NumTypeInt16:
		return math.MinInt16, math.MaxInt16
	case NumTypeUint32:
		return 0, math.MaxUint32
	case NumTypeInt32:
		return math.MinInt32, math.MaxInt32
	case NumTypeFloat32, NumTypeFloat64:
		return math.MinInt64, math.MaxFloat64
	}

	return 0, 0
}

// IntoType will convert a `NumType` to a `jtd.Type`.
func (n NumType) IntoType() jtd.Type {
	switch n {
	case NumTypeUint8:
		return jtd.TypeUint8
	case NumTypeInt8:
		return jtd.TypeInt8
	case NumTypeUint16:
		return jtd.TypeUint16
	case NumTypeInt16:
		return jtd.TypeInt16
	case NumTypeUint32:
		return jtd.TypeUint32
	case NumTypeInt32:
		return jtd.TypeInt32
	case NumTypeFloat32:
		return jtd.TypeFloat32
	case NumTypeFloat64:
		return jtd.TypeFloat64
	}

	return jtd.TypeUint8
}

// InferredNumber represents the state for a column that is a number. It holds
// the seen maximum and minimum value together with information about if all
// seen numbers are integers.
type InferredNumber struct {
	Min       float64
	Max       float64
	IsInteger bool
}

// NewNumber will return a new `InferredNumber`.
func NewNumber() *InferredNumber {
	return &InferredNumber{
		IsInteger: true,
	}
}

// Infer will infer a value, updating the state for the `InferredNumber`.
func (i *InferredNumber) Infer(n float64) *InferredNumber {
	return &InferredNumber{
		Min:       math.Min(i.Min, n),
		Max:       math.Max(i.Max, n),
		IsInteger: i.IsInteger && float64(int(n)) == n,
	}
}

// InfoType will convert an `InferredNumber` to a `jtd.Type`.
func (i *InferredNumber) IntoType(defaultType NumType) jtd.Type {
	if i.ContainedBy(defaultType) {
		return defaultType.IntoType()
	}

	numTypes := []NumType{
		NumTypeUint8,
		NumTypeInt8,
		NumTypeUint16,
		NumTypeInt16,
		NumTypeUint32,
		NumTypeInt32,
	}

	for _, v := range numTypes {
		if i.ContainedBy(v) {
			return v.IntoType()
		}
	}

	return jtd.TypeFloat64
}

// ContainedBy checks if an inferred number column can be contained within the
// passed `NumType`, meaning it is above the minimum and below the maximum value
// for the number type.
func (i *InferredNumber) ContainedBy(nt NumType) bool {
	if !i.IsInteger && !nt.IsFloat() {
		return false
	}

	min, max := nt.AsRange()

	return min <= i.Min && max >= i.Max
}
