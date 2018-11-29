package yassembly

import (
	"fmt"
	"io"
	"strings"
)

func Merge(base io.Reader, imports []io.Reader, output io.Writer) error {
	_, err := io.Copy(output, base)
	if err != nil {
		return fmt.Errorf("Copy of base errored with: %v", err)
	}

	for _, i := range imports {
		_, err := io.Copy(output, strings.NewReader("\n"))
		if err != nil {
			return fmt.Errorf("Copy of newline errored with: %v", err)
		}
		_, err = io.Copy(output, i)
		if err != nil {
			return fmt.Errorf("Copy of import errored with: %v", err)
		}
	}
	return nil
}
