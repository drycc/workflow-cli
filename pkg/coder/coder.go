package coder

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
)

// tomlType builds a mirror type of t at runtime: fields keep the SDK struct
// definition order, keys are named after the json tags, and both json/toml
// tags are attached so the type can be JSON-decoded and TOML-encoded.
// Struct types without exported json-tagged fields (e.g. time.Time) are
// treated as opaque values and returned unchanged.
func tomlType(t reflect.Type) reflect.Type {
	switch t.Kind() {
	case reflect.Pointer:
		return reflect.PointerTo(tomlType(t.Elem()))
	case reflect.Slice, reflect.Array:
		return reflect.SliceOf(tomlType(t.Elem()))
	case reflect.Map:
		return t
	case reflect.Struct:
		fields := make([]reflect.StructField, 0, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if !f.IsExported() {
				continue
			}
			jsonTag := f.Tag.Get("json")
			name := strings.Split(jsonTag, ",")[0]
			if name == "-" {
				continue
			}
			if name == "" {
				name = f.Name
			}
			fieldType := tomlType(f.Type)
			tomlTag := `json:"` + jsonTag + `" toml:"` + name + `"`
			if strings.Contains(jsonTag, "omitempty") {
				tomlTag = `json:"` + jsonTag + `" toml:"` + name + `,omitempty"`
				// BurntSushi/toml 编码不支持 omitempty；标量提升为指针，
				// JSON 中缺省时保持 nil，编码器即跳过
				switch fieldType.Kind() {
				case reflect.String, reflect.Bool,
					reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
					reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
					reflect.Float32, reflect.Float64:
					fieldType = reflect.PointerTo(fieldType)
				}
			}
			fields = append(fields, reflect.StructField{
				Name: f.Name,
				Type: fieldType,
				Tag:  reflect.StructTag(tomlTag),
			})
		}
		if len(fields) == 0 {
			return t
		}
		return reflect.StructOf(fields)
	default:
		return t
	}
}

// MarshalTOML encodes an SDK info struct as TOML, preserving the struct's
// field definition order. json.Number is kept so integers stay integers.
func MarshalTOML(info any) ([]byte, error) {
	data, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	outPtr := reflect.New(tomlType(reflect.TypeOf(info)))
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(outPtr.Interface()); err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(outPtr.Interface()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalTOML decodes TOML data into an SDK request struct via JSON, so
// json tags on the SDK structs apply without any duplicated definitions.
func UnmarshalTOML(data []byte, out any) error {
	var m map[string]any
	if err := toml.Unmarshal(data, &m); err != nil {
		return err
	}
	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, out)
}
