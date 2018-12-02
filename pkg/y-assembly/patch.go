package yassembly

import (
	"bytes"
	"fmt"
	"io"

	yamlpatch "github.com/krishicks/yaml-patch"
)

func Patch(base io.Reader, patches []io.Reader, output io.Writer) error {
	if len(patches) == 0 {
		_, err := io.Copy(output, base)
		if err != nil {
			return fmt.Errorf("Copy of base errored with: %v", err)
		}
		return nil
	}

	baseBuffer := new(bytes.Buffer)
	baseBuffer.ReadFrom(base)

	for _, patchReader := range patches {
		ops := new(bytes.Buffer)
		ops.ReadFrom(patchReader)
		patch, err := yamlpatch.DecodePatch(ops.Bytes())
		if err != nil {
			return fmt.Errorf("decoding patch failed: %s", err)
		}

		modifiedBytes, err := patch.Apply(baseBuffer.Bytes())
		if err != nil {
			return fmt.Errorf("applying patch failed: %s", err)
		}

		modifiedBuffer := bytes.NewBuffer(modifiedBytes)
		_, err = io.Copy(output, modifiedBuffer)
		if err != nil {
			return fmt.Errorf("Copy of patch errored with: %v", err)
		}
	}
	return nil
}
