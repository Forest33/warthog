package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

const (
	defaultConfigFile = "warthog.json"
	tagDefault        = "default"
)

func init() {
	cfgPath, ok := os.LookupEnv("WARTHOG_CONFIG")
	if !ok {
		ex, err := os.Executable()
		if err != nil {
			log.Fatal(err.Error())
		}
		cfgPath = filepath.Join(filepath.Dir(ex), defaultConfigFile)
	}

	data, err := os.ReadFile(cfgPath)
	if err == nil {
		err = json.Unmarshal(data, cfg)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	if err := parse(cfg); err != nil {
		log.Fatal(err.Error())
	}
}

func parse(target interface{}) error {
	ref := reflect.Indirect(reflect.ValueOf(target))
	for i := 0; i < ref.Type().NumField(); i++ {
		structField := ref.Type().Field(i)
		fieldValue := ref.Field(i)

		defaultTagValue, defaultTagExists := structField.Tag.Lookup(tagDefault)

		if defaultTagExists {
			if err := setValue(structField, &fieldValue, defaultTagValue); err != nil {
				return err
			}
			continue
		}

		if structField.Type.Kind() != reflect.Ptr {
			return fmt.Errorf("required configuration parameter is not specified - %s.%s", ref.Type().Name(), structField.Name)
		}

		if err := setValue(structField, &fieldValue, ""); err != nil {
			return err
		}
	}

	return nil
}

func setValue(structField reflect.StructField, field *reflect.Value, value string) error {
	switch structField.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(value, 10, int(structField.Type.Size()*8))
		if err != nil {
			return err
		}
		field.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(value, 10, int(structField.Type.Size()*8))
		if err != nil {
			return err
		}
		field.SetUint(v)
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(value, int(structField.Type.Size()*8))
		if err != nil {
			return err
		}
		field.SetFloat(v)
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		field.SetBool(strings.ToLower(value) == "true")
	case reflect.Ptr:
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		return parse(field.Interface())
	case reflect.Slice:
		if len(value) > 0 {
			values := strings.Split(value, ",")
			sl := reflect.MakeSlice(field.Type(), len(values), len(values))
			for i, val := range values {
				sl.Index(i).Set(reflect.ValueOf(val))
			}
			field.Set(sl)
		}
	}
	return nil
}
