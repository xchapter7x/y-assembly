package yup_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/xchapter7x/yup/pkg/yup"
)

func TestConfigParse(t *testing.T) {
	controlPathBase := "testpath"
	singleImport := strings.NewReader(`---
- version: 1
  base: "somefile.yml"
  output: "somenewfile.yml"
  imports:
  - "` + controlPathBase + `/local/file"`)

	multipleImports := strings.NewReader(`---
- version: 1
  base: "somefile.yml"
  output: "somenewfile.yml"
  imports:
  - "` + controlPathBase + `/local/file"
  - "` + controlPathBase + `/other/local/file"`)

	invalidYupfile := strings.NewReader(`!*#) ../:idddkgh`)
	noImportsYupfile := strings.NewReader(`version: 1`)
	emptyYupfile := strings.NewReader(``)

	noOutputYupfile := strings.NewReader(`---
- version: 1
  output: "somenewfile.yml"
  imports:
  - "` + controlPathBase + `/local/file"`)

	noBaseYupfile := strings.NewReader(`---
- version: 1
  output: "somenewfile.yml"
  imports:
  - "` + controlPathBase + `/local/file"`)

	t.Run("parse success", func(t *testing.T) {
		for _, table := range []struct {
			testName   string
			fileReader io.Reader
		}{
			{"defined single local import", singleImport},
			{"defined multiple local import", multipleImports},
		} {
			t.Run(table.testName, func(t *testing.T) {
				configs, err := yup.ConfigParse(table.fileReader)
				if len(configs) == 0 {
					t.Error("Expected to have configs")
				}
				for _, config := range configs {
					for _, importPath := range config.Imports {
						if !strings.HasPrefix(importPath, controlPathBase) {
							t.Errorf("expected path prefix to be %s in path %s", controlPathBase, importPath)
						}
					}
				}

				if configs == nil {
					t.Error("Expected to have config got nil")
				}

				if err != nil {
					t.Errorf("Error response: %v", err)
				}

			})
		}
	})

	t.Run("parse failures", func(t *testing.T) {
		for _, table := range []struct {
			testName   string
			fileReader io.Reader
		}{
			{"invalid yupfile", invalidYupfile},
			{"no imports", noImportsYupfile},
			{"no output yaml", noOutputYupfile},
			{"no base yaml", noBaseYupfile},
			{"empty yupfile", emptyYupfile},
		} {
			t.Run(table.testName, func(t *testing.T) {
				configs, err := yup.ConfigParse(table.fileReader)
				fmt.Println(err)
				if configs != nil {
					t.Errorf("Expected to have empty config: %v", configs)
				}
				if err == nil {
					t.Error("Expected error but got nil")
				}
			})
		}
	})
}
