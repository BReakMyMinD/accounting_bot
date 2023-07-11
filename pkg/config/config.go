package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
)

type configReader struct {
	config map[string]interface{}
}

func NewConfigReader(path string) (*configReader, error) {
	jsonConfig, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonConfig.Close()
	byteConfig, err := io.ReadAll(jsonConfig)
	if err != nil {
		return nil, err
	}
	var reader configReader
	err = json.Unmarshal(byteConfig, &reader.config)
	if err != nil {
		return nil, err
	} else {
		return &reader, nil
	}
}

func (configReader *configReader) GetParameter(name string, valuePtr interface{}) error {
	if configValue, ok := configReader.config[name]; ok && configValue != nil {
		configValuePtrType := reflect.PointerTo(reflect.TypeOf(configValue))
		valuePtrType := reflect.TypeOf(valuePtr)
		if configValuePtrType != valuePtrType {
			return fmt.Errorf("pointer to parameter %q has type %s, passed %T", name, configValuePtrType.String(), valuePtr)
		} else {
			reflect.ValueOf(valuePtr).Elem().Set(reflect.ValueOf(configValue))
			return nil
		}
	} else {
		return fmt.Errorf("parameter %q is not found", name)
	}
}