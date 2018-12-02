package yassembly_test

import (
	"bytes"
	"io"
	"testing"

	yassembly "github.com/xchapter7x/y-assembly/pkg/y-assembly"
)

func TestPatch(t *testing.T) {
	sampleDoc := bytes.NewBuffer([]byte(`---
foo: bar
baz:
  quux: grault
`))

	sampleOps := bytes.NewBuffer([]byte(`---
- op: add
  path: /baz/waldo
  value: fred
`))
	t.Run("should apply patch", func(t *testing.T) {
		testOutput := bytes.NewBuffer([]byte(``))
		yassembly.Patch(sampleDoc, []io.Reader{sampleOps}, testOutput)
		controlOutput := bytes.NewBuffer([]byte(`baz:
  quux: grault
  waldo: fred
foo: bar
`))
		if controlOutput.String() != testOutput.String() {
			t.Errorf("%s != %s",
				controlOutput.String(),
				testOutput.String(),
			)
		}
	})
}
