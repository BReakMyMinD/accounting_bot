package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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

func (configReader *configReader) GetString(name string) (string, error) {
	if value, ok := configReader.config[name]; ok {
		if stringValue, ok := value.(string); ok {
			return stringValue, nil
		} else {
			return "", fmt.Errorf("parameter %q is not a string", name)
		}
	} else {
		return "", fmt.Errorf("parameter %q is not found", name)
	}
}