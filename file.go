package main

import (
	"encoding/json"
	"strings"

	toml "github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func blobToObject(fileType string, blob []byte) (object interface{}, err error) {
	switch fileType {
	case "TOML":
		t, e := toml.LoadBytes(blob)
		if e != nil {
			return nil, e
		}
		object = t.ToMap()
	case "YAML":
		m := map[string]interface{}{}
		e := yaml.Unmarshal(blob, &m)
		if e != nil {
			return nil, e
		}
		object = m
	case "JSON":
		m := map[string]interface{}{}
		e := json.Unmarshal(blob, &m)
		if e != nil {
			return nil, e
		}
		object = m
	case "TEXT":
		object = string(blob)
	case "STRING":
		object = strings.TrimSpace(string(blob))
	case "RAW":
		object = blob
	default:
		return nil, errors.Errorf("%s is unknown type", fileType)
	}

	return object, nil
}

func getFileObject(s string, defaultFileType string) (fileType string, filePath string, object interface{}, err error) {
	f := fileArg(s)
	fileType, filePath = f.parse()

	blob, err := f.blob()
	if err != nil {
		return fileType, filePath, nil, err
	}

	if fileType == "" {
		fileType = defaultFileType
	}

	object, err = blobToObject(fileType, blob)
	if err != nil {
		return fileType, filePath, nil, err
	}

	return fileType, filePath, object, err
}
