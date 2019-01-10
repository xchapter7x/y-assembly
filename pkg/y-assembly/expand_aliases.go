package yassembly

import (
	"bytes"
	"fmt"
	"io"

	yaml "gopkg.in/yaml.v2"
)

func ExpandAliases(input io.Reader, output io.Writer) error {
	inputBuffer := new(bytes.Buffer)
	inputBuffer.ReadFrom(input)
	var expandedYaml interface{}
	err := yaml.Unmarshal(inputBuffer.Bytes(), &expandedYaml)
	if err != nil {
		return fmt.Errorf("couldnt unmarshal yaml: %v", err)
	}

	expandedBytes, err := yaml.Marshal(expandedYaml)
	if err != nil {
		return fmt.Errorf("failed marshaling yaml: %v", err)
	}

	expandedBuffer := bytes.NewBuffer(expandedBytes)
	_, err = io.Copy(output, expandedBuffer)
	if err != nil {
		return fmt.Errorf("Copy of expanded aliases errored with: %v", err)
	}
	return nil
}
