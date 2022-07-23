package cfg

import (
	"errors"
	"os"
	"reflect"
	"strconv"
)

// tag in the structure by which values in global env will be searched for
const tagName = "env"

// LoadFromEnv Initialize configuration
func LoadFromEnv(config any) error {
	val := reflect.ValueOf(config)
	err := unmarshalConfig(val)
	if err != nil {
		return err
	}
	return nil
}

// unmarshalConfig load data into configuration fields by tag tagName
func unmarshalConfig(v reflect.Value) error {
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			subStruct := v.Field(i).Addr().Elem()
			err := unmarshalConfig(subStruct)
			if err != nil {
				return err
			}
			continue
		}
		if tag := field.Tag.Get(tagName); tag == "-" || tag == "" {
			continue
		}
		if err := setValue(v, field); err != nil {
			return err
		}
	}
	return nil
}

// setValue assigns the field.Name field a value from the global ENV
func setValue(ps reflect.Value, field reflect.StructField) error {
	tag := field.Tag.Get("env")
	f := ps.FieldByName(field.Name)
	if f.IsValid() {
		if f.Kind() == reflect.Pointer {
			f.Set(reflect.New(f.Type().Elem()))
			f = f.Elem()
		}
		s, err := setEnvConfig(tag, f.Kind())
		if err != nil {
			if errors.As(err, &GetEnvError{}) {
				return nil
			}
			return err
		}
		f.Set(s)
	}
	return nil
}

// setEnvConfig returns reflect.Value reduced to the field type of the structure
func setEnvConfig(name string, kind reflect.Kind) (reflect.Value, error) {
	if value := os.Getenv(name); len(value) > 0 {
		switch kind {
		case reflect.String:
			return reflect.ValueOf(value), nil
		case reflect.Int:
			i, err := strconv.Atoi(value)
			return reflect.ValueOf(i), err
		case reflect.Bool:
			b, err := strconv.ParseBool(value)
			return reflect.ValueOf(b), err
		default:
			return reflect.Value{}, errors.New("no converter to the right data type: " + kind.String() + " was found ")
		}
	}
	return reflect.Value{}, GetEnvError{field: name}
}
