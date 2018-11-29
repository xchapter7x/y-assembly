package yassembly

import (
	"bytes"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Version int
	Base    string
	Output  string
	Imports []string
}

func ConfigParse(f io.Reader) ([]Config, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(f)
	c := make([]Config, 0)
	b := buf.Bytes()
	if len(b) == 0 {
		return nil, fmt.Errorf("The provided reader was empty")
	}

	err := yaml.Unmarshal(buf.Bytes(), &c)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling: %v \n%s", err, string(b))
	}

	if len(c) == 0 {
		return nil, fmt.Errorf("No configs defined: %v", string(b))
	}

	for _, v := range c {
		if v.Base == "" {
			return nil, fmt.Errorf("Base must be defined for: %v", v)
		}
		if v.Output == "" {
			return nil, fmt.Errorf("Output must be defined for: %v", v)
		}
		if len(v.Imports) == 0 {
			return nil, fmt.Errorf("No imports configured on: %v", v)
		}
	}
	return c, err
}
