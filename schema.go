package jtdinfer

import jtd "github.com/jsontypedef/json-typedef-go"

// Schema represents the JTD schema that will get inferred. It's a hard copy of
// the upstream type since we want and need to be able to omit empty fields.
// Ref: https://github.com/jsontypedef/json-typedef-go/issues/7
type Schema struct {
	Definitions          map[string]Schema      `json:"definitions,omitempty"`
	Metadata             map[string]interface{} `json:"metadata,omitempty"`
	Nullable             bool                   `json:"nullable,omitempty"`
	Ref                  *string                `json:"ref,omitempty"`
	Type                 jtd.Type               `json:"type,omitempty"`
	Enum                 []string               `json:"enum,omitempty"`
	Elements             *Schema                `json:"elements,omitempty"`
	Properties           map[string]Schema      `json:"properties,omitempty"`
	OptionalProperties   map[string]Schema      `json:"optionalProperties,omitempty"`
	AdditionalProperties bool                   `json:"additionalProperties,omitempty"`
	Values               *Schema                `json:"values,omitempty"`
	Discriminator        string                 `json:"discriminator,omitempty"`
	Mapping              map[string]Schema      `json:"mapping,omitempty"`
}
