package yassembly_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/xchapter7x/y-assembly/pkg/y-assembly"
)

func TestCombine(t *testing.T) {
	t.Run("valid base yaml and imports", func(t *testing.T) {
		controlLineOne := `---`
		controlLineTwo := `a:a`
		controlLineThree := `b:b`
		baseFixture := strings.NewReader(controlLineOne)
		importsFixture := []io.Reader{
			strings.NewReader(controlLineTwo),
			strings.NewReader(controlLineThree),
		}
		controlOutput := fmt.Sprintf("%s\n%s\n%s", controlLineOne, controlLineTwo, controlLineThree)

		outputFixture := new(strings.Builder)
		err := yassembly.Combine(baseFixture, importsFixture, outputFixture)

		if err != nil {
			t.Errorf("did not expect error: %v", err)
		}

		if outputFixture.String() != controlOutput {
			t.Errorf("string on output incorrect:\n%s", outputFixture.String())
		}
	})
}
