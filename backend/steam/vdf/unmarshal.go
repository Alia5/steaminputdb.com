package vdf

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type UnmarshalTypeError struct {
	Value  string
	Type   reflect.Type
	Struct string
	Field  string
}

func (e *UnmarshalTypeError) Error() string {
	if e.Struct != "" || e.Field != "" {
		return "cannot unmarshal " + e.Value + " into Go struct field " + e.Struct + "." + e.Field + " of type " + e.Type.String()
	}
	return "cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}

type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "Unmarshal(nil)"
	}
	if e.Type.Kind() != reflect.Pointer {
		return "Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "Unmarshal(nil " + e.Type.String() + ")"
}

func Unmarshal(data string, v any) error {
	m, err := Parse(data)
	if err != nil {
		return err
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}
	return decodeAny(m, rv.Elem())
}

func decodeAny(val any, rv reflect.Value) error {
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}

	switch v := val.(type) {
	case map[string]any:
		return decodeMapAny(v, rv)
	case []any:
		return decodeSliceAny(v, rv)
	case string:
		return decodeLeaf(v, rv)
	case int64:
		return decodeInt(v, rv)
	case float64:
		return decodeFloat(v, rv)
	case bool:
		if v {
			return decodeLeaf("true", rv)
		}
		return decodeLeaf("false", rv)
	default:
		return &UnmarshalTypeError{Value: fmt.Sprintf("%T", val), Type: rv.Type()}
	}
}

func decodeMapAny(m map[string]any, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Struct:
		return decodeStruct(m, rv)
	case reflect.Map:
		return decodeMapValue(m, rv)
	case reflect.Interface:
		rv.Set(reflect.ValueOf(m))
		return nil
	default:
		return &UnmarshalTypeError{Value: "map", Type: rv.Type()}
	}
}

func decodeSliceAny(arr []any, rv reflect.Value) error {
	if rv.Kind() == reflect.Interface {
		rv.Set(reflect.ValueOf(arr))
		return nil
	}
	if rv.Kind() != reflect.Slice {
		return &UnmarshalTypeError{Value: "array", Type: rv.Type()}
	}
	for _, elem := range arr {
		ev := reflect.New(rv.Type().Elem()).Elem()
		if err := decodeAny(elem, ev); err != nil {
			return err
		}
		rv.Set(reflect.Append(rv, ev))
	}
	return nil
}

func decodeStruct(m map[string]any, rv reflect.Value) error {
	rt := rv.Type()
	for i := range rt.NumField() {
		field := rt.Field(i)
		if !field.IsExported() {
			continue
		}
		key := fieldKey(field)
		if key == "-" {
			continue
		}
		val, ok := m[key]
		if !ok {
			continue
		}
		fv := rv.Field(i)
		if err := decodeAny(val, fv); err != nil {
			if te, ok := err.(*UnmarshalTypeError); ok {
				te.Struct = rt.Name()
				te.Field = field.Name
			}
			return err
		}
	}
	return nil
}

func decodeLeaf(val string, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.String:
		rv.SetString(val)
		return nil
	case reflect.Bool:
		rv.SetBool(val == "1" || strings.EqualFold(val, "true"))
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		rv.SetInt(n)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		rv.SetUint(n)
		return nil
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		rv.SetFloat(n)
		return nil
	case reflect.Interface:
		rv.Set(reflect.ValueOf(val))
		return nil
	default:
		return &UnmarshalTypeError{Value: "string", Type: rv.Type()}
	}
}

func decodeMapValue(m map[string]any, rv reflect.Value) error {
	rt := rv.Type()
	if rv.IsNil() {
		rv.Set(reflect.MakeMap(rt))
	}
	keyType := rt.Key()
	if keyType.Kind() != reflect.String {
		return &UnmarshalTypeError{Value: "map", Type: rt}
	}
	elemType := rt.Elem()
	for k, val := range m {
		ev := reflect.New(elemType).Elem()
		if err := decodeAny(val, ev); err != nil {
			return err
		}
		rv.SetMapIndex(reflect.ValueOf(k), ev)
	}
	return nil
}

func fieldKey(f reflect.StructField) string {
	tag := f.Tag.Get("json")
	if tag == "" {
		tag = f.Tag.Get("vdf")
	}
	if tag == "" {
		return strings.ToLower(f.Name)
	}
	key, _, _ := strings.Cut(tag, ",")
	if key == "" {
		return strings.ToLower(f.Name)
	}
	return key
}

func decodeInt(v int64, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rv.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		rv.SetUint(uint64(v))
	case reflect.Float32, reflect.Float64:
		rv.SetFloat(float64(v))
	case reflect.Bool:
		rv.SetBool(v != 0)
	case reflect.String:
		rv.SetString(strconv.FormatInt(v, 10))
	case reflect.Interface:
		rv.Set(reflect.ValueOf(v))
	default:
		return &UnmarshalTypeError{Value: "number", Type: rv.Type()}
	}
	return nil
}

func decodeFloat(v float64, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Float32, reflect.Float64:
		rv.SetFloat(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rv.SetInt(int64(v))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		rv.SetUint(uint64(v))
	case reflect.String:
		rv.SetString(strconv.FormatFloat(v, 'f', -1, 64))
	case reflect.Interface:
		rv.Set(reflect.ValueOf(v))
	default:
		return &UnmarshalTypeError{Value: "number", Type: rv.Type()}
	}
	return nil
}
