package jtdinfer

import (
	"fmt"
	"math"
	"reflect"
	"time"

	jtd "github.com/jsontypedef/json-typedef-go"
)

type SchemaType int8

const (
	SchemaTypeUnknown SchemaType = iota
	SchemaTypeAny
	SchemaTypeBoolean
	SchemaTypeNumber
	SchemaTypeString
	SchemaTypeTimestmap
	SchemaTypeEnum
	SchemaTypeArray
	SchemaTypeProperties
	SchemaTypeValues
	SchemaTypeDiscriminator
	SchemaTypeNullable
)

type Properties struct {
	Required map[string]*InferredSchema
	Optional map[string]*InferredSchema
}

type Discriminator struct {
	Discriminator string
	Mapping       map[string]*InferredSchema
}

// InferredSchema is the schema while being inferred that holds information
// about fields.
type InferredSchema struct {
	SchemaType    SchemaType
	Number        *InferredNumber
	Enum          map[string]struct{}
	Array         *InferredSchema
	Properties    Properties
	Values        *InferredSchema
	Discriminator Discriminator
	Nullable      *InferredSchema
}

func NewInferredSchema() *InferredSchema {
	return &InferredSchema{}
}

func NewInferredSchemaWithType(t SchemaType) *InferredSchema {
	return &InferredSchema{SchemaType: t}
}

// Infer will infer the schema by trying to mimic the way it's implemented in
// the Rust implementation at
// https://github.com/jsontypedef/json-typedef-infer/blob/master/src/inferred_schema.rs.
// Since we don't have enums of this kind in Go we're using a struct with
// pointers to a schema instead of wrapping the enums.
func (i *InferredSchema) Infer(value any) *InferredSchema {
	if value == nil {
		return &InferredSchema{
			SchemaType: SchemaTypeNullable,
			Nullable:   i,
		}
	}

	if i.SchemaType == SchemaTypeNullable {
		return &InferredSchema{
			SchemaType: SchemaTypeNullable,
			Nullable:   i,
		}
	}

	if _, ok := value.(bool); ok && i.SchemaType == SchemaTypeUnknown {
		return &InferredSchema{SchemaType: SchemaTypeBoolean}
	}

	// In Go all numbers from unmarshalling JSON will be represented as float64
	// so this is the only type we need.
	if v, ok := value.(float64); ok && i.SchemaType == SchemaTypeUnknown {
		return &InferredSchema{
			SchemaType: SchemaTypeNumber,
			Number:     NewNumber().Infer(v),
		}
	}

	if v, ok := value.(string); ok && i.SchemaType == SchemaTypeUnknown {
		// TODO: Enums
		if _, err := time.Parse(time.RFC3339, v); err == nil {
			return &InferredSchema{SchemaType: SchemaTypeTimestmap}
		}

		return &InferredSchema{SchemaType: SchemaTypeString}
	}

	if i.SchemaType == SchemaTypeUnknown && reflect.TypeOf(value).Kind() == reflect.Slice {
		s := reflect.ValueOf(value)

		subInfer := &InferredSchema{}
		for i := 0; i < s.Len(); i++ {
			subInfer = subInfer.Infer(s.Index(i).Interface())
		}

		return &InferredSchema{
			SchemaType: SchemaTypeArray,
			Array:      subInfer,
		}
	}

	if m, ok := value.(map[string]any); ok && i.SchemaType == SchemaTypeUnknown {
		// TODO: Hints
		// TODO: Discriminator
		properties := make(map[string]*InferredSchema, 0)
		for k, v := range m {
			properties[k] = NewInferredSchema().Infer(v)
		}

		return &InferredSchema{
			SchemaType: SchemaTypeProperties,
			Properties: Properties{
				Required: properties,
			},
		}
	}

	if i.SchemaType == SchemaTypeAny {
		return &InferredSchema{SchemaType: SchemaTypeAny}
	}

	if _, ok := value.(bool); ok && i.SchemaType == SchemaTypeBoolean {
		return &InferredSchema{SchemaType: SchemaTypeBoolean}
	}

	if i.SchemaType == SchemaTypeBoolean {
		return &InferredSchema{SchemaType: SchemaTypeAny}
	}

	if v, ok := value.(float64); ok && i.SchemaType == SchemaTypeNumber {
		return &InferredSchema{
			SchemaType: SchemaTypeNumber,
			Number:     i.Number.Infer(v),
		}
	}

	if i.SchemaType == SchemaTypeNumber {
		return &InferredSchema{SchemaType: SchemaTypeAny}
	}

	if v, ok := value.(string); ok && i.SchemaType == SchemaTypeTimestmap {
		if _, err := time.Parse(time.RFC3339, v); err == nil {
			return &InferredSchema{SchemaType: SchemaTypeTimestmap}
		}

		return &InferredSchema{SchemaType: SchemaTypeString}
	}

	if i.SchemaType == SchemaTypeTimestmap {
		return &InferredSchema{SchemaType: SchemaTypeAny}
	}

	if _, ok := value.(string); ok && i.SchemaType == SchemaTypeString {
		return &InferredSchema{SchemaType: SchemaTypeString}
	}

	if i.SchemaType == SchemaTypeString {
		return &InferredSchema{SchemaType: SchemaTypeAny}
	}

	if v, ok := value.(string); ok && i.SchemaType == SchemaTypeEnum {
		i.Enum[v] = struct{}{}
	}

	if i.SchemaType == SchemaTypeEnum {
		return &InferredSchema{SchemaType: SchemaTypeAny}
	}

	if i.SchemaType == SchemaTypeArray && reflect.TypeOf(value).Kind() == reflect.Slice {
		// TODO: Hints
		s := reflect.ValueOf(value)

		subInfer := i.Array
		for i := 0; i < s.Len(); i++ {
			subInfer = subInfer.Infer(s.Index(i).Interface())
		}

		return &InferredSchema{
			SchemaType: SchemaTypeArray,
			Array:      subInfer,
		}
	}

	if m, ok := value.(map[string]any); ok && i.SchemaType == SchemaTypeProperties {
		ensureMap := func(m map[string]*InferredSchema) map[string]*InferredSchema {
			if m != nil {
				return m
			}

			return make(map[string]*InferredSchema, 0)
		}

		missingKeys := []string{}

		for k := range i.Properties.Required {
			if _, ok := m[k]; !ok {
				missingKeys = append(missingKeys, k)
			}
		}

		for _, k := range missingKeys {
			subInfer := i.Properties.Required[k]
			delete(i.Properties.Required, k)

			i.Properties.Optional = ensureMap(i.Properties.Optional)
			i.Properties.Optional[k] = subInfer
		}

		for k, v := range m {
			if subInfer, ok := i.Properties.Required[k]; ok {
				i.Properties.Required[k] = subInfer.Infer(v)
			} else if subInfer, ok := i.Properties.Optional[k]; ok {
				i.Properties.Optional = ensureMap(i.Properties.Optional)
				i.Properties.Optional[k] = subInfer.Infer(v)
			} else {
				i.Properties.Optional = ensureMap(i.Properties.Optional)
				i.Properties.Optional[k] = NewInferredSchema().Infer(v)
			}
		}

		return i
	}

	if i.SchemaType == SchemaTypeProperties {
		return &InferredSchema{SchemaType: SchemaTypeAny}
	}

	if m, ok := value.(map[string]any); ok && i.SchemaType == SchemaTypeValues {
		// TODO: Hints
		subInfer := i.Values
		for _, v := range m {
			subInfer = NewInferredSchema().Infer(v)
		}

		return &InferredSchema{
			SchemaType: SchemaTypeValues,
			Values:     subInfer,
		}
	}

	if i.SchemaType == SchemaTypeValues {
		return &InferredSchema{SchemaType: SchemaTypeAny}
	}

	if m, ok := value.(map[string]any); ok && i.SchemaType == SchemaTypeDiscriminator {
		mappingKey, ok := m[i.Discriminator.Discriminator].(string)
		if !ok {
			return &InferredSchema{SchemaType: SchemaTypeAny}
		}

		if _, ok := i.Discriminator.Mapping[mappingKey]; !ok {
			i.Discriminator.Mapping[mappingKey] = NewInferredSchema()
		}

		i.Discriminator.Mapping[mappingKey] = i.Discriminator.Mapping[mappingKey].Infer(m)
	}

	if i.SchemaType == SchemaTypeDiscriminator {
		return &InferredSchema{SchemaType: SchemaTypeAny}
	}

	panic(fmt.Sprintf("%T: %T (%v)", i.SchemaType, value, value))
}

func (i *InferredSchema) IntoSchema() Schema {
	switch i.SchemaType {
	case SchemaTypeUnknown, SchemaTypeAny:
		return Schema{}

	case SchemaTypeBoolean:
		return Schema{Type: jtd.TypeBoolean}
	case SchemaTypeNumber:
		return Schema{
			Type: i.Number.IntoType(minMax{
				typ: jtd.TypeUint8,
				min: 0,
				max: math.MaxUint8,
			}),
		}
	case SchemaTypeString:
		return Schema{Type: jtd.TypeString}
	case SchemaTypeTimestmap:
		return Schema{Type: jtd.TypeTimestamp}
	case SchemaTypeEnum:
		enum := make([]string, 0, len(i.Enum))
		for k := range i.Enum {
			enum = append(enum, k)
		}

		return Schema{Enum: enum}
	case SchemaTypeArray:
		elements := i.Array.IntoSchema()
		return Schema{Elements: &elements}
	case SchemaTypeProperties:
		var (
			required map[string]Schema
			optional map[string]Schema
		)

		if i.Properties.Required != nil {
			required = make(map[string]Schema, len(i.Properties.Required))

			for k, v := range i.Properties.Required {
				required[k] = v.IntoSchema()
			}
		}

		if i.Properties.Optional != nil {
			optional = make(map[string]Schema, len(i.Properties.Optional))

			for k, v := range i.Properties.Optional {
				optional[k] = v.IntoSchema()
			}
		}

		return Schema{
			Properties:         required,
			OptionalProperties: optional,
		}
	case SchemaTypeValues:
		values := i.Values.IntoSchema()
		return Schema{Values: &values}
	case SchemaTypeDiscriminator:
		// TODO: Add support for discriminator

	case SchemaTypeNullable:
		schema := i.Nullable.IntoSchema()
		schema.Nullable = true

		return schema
	}

	return Schema{}
}
