# JSON typedef infer

[![Go Reference](https://pkg.go.dev/badge/github.com/bombsimon/jtd-infer-go.svg)](https://pkg.go.dev/github.com/bombsimon/jtd-infer-go)
[![Build and test](https://github.com/bombsimon/jtd-infer-go/actions/workflows/go.yaml/badge.svg)](https://github.com/bombsimon/jtd-infer-go/actions/workflows/go.yaml)

This is a port of [`json-typedef-infer`][jtd-infer] for Go. The reason for
porting this is that I was in need of JTD inference from code and not as a CLI
tool.

For more information about JSON Typedef and its RFC and how to use different
kind of hints see [`json-typedef-infer`][jtd-infer].

## Usage

See [examples] directory for runnable examples and how to infer JTD. If you want
to infer a primitive Go type such as a `string` or `map[string]any`, you can use
`Infer`:

```go
schema := NewInferrer(WithoutHints()).
    Infer("my-string").
    IntoSchema()
// {
//   "type": "string"
// }
```

```go
schema := NewInferrer(WithoutHints()).
    Infer(map[string]any{
        "age":  52,
        "name": "Joe",
    }).
    IntoSchema()
// {
//   "properties": {
//     "age": {
//       "type": "uint8"
//     },
//     "name": {
//       "type": "string"
//     }
//   }
// }
```

If you have one or more rows of JSON objects or lists as strings you can pass
them to the shorthand function `InferStrings`. This will create an `Inferrer`
and call `Infer` repeatedly on each row after deserializing it to a Go map. This
is the recommended way to infer the schema even if you just have a single line.

```go
rows := []string{
    `{"name":"Joe", "age": 52, "something_optional": true, "something_nullable": 1.1}`,
    `{"name":"Jane", "age": 48, "something_nullable": null}`,
}
schema := InferStrings(rows, WithoutHints()).IntoSchema()
// {
//   "properties": {
//     "age": {
//       "type": "uint8"
//     },
//     "name": {
//       "type": "string"
//     },
//     "something_nullable": {
//       "nullable": true,
//       "type": "float64"
//     }
//   },
//   "optionalProperties": {
//     "something_optional": {
//       "type": "boolean"
//     }
//   }
// }
```

[jtd-infer]: https://github.com/jsontypedef/json-typedef-infer/
[examples]: examples
