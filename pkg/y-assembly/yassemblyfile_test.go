package yassembly_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/xchapter7x/y-assembly/pkg/y-assembly"
)

func TestConfigParse(t *testing.T) {
	controlPathBase := "testpath"
	singlePatch := strings.NewReader(`---
- version: 1
  base: "somefile.yml"
  output: "somenewfile.yml"
  imports:
  - "` + controlPathBase + `/local/file"
  patches:
  - "` + controlPathBase + `/local/patch"`)

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

	invalidYassemblyfile := strings.NewReader(`!*#) ../:idddkgh`)
	noImportsYassemblyfile := strings.NewReader(`version: 1`)
	emptyYassemblyfile := strings.NewReader(``)

	noOutputYassemblyfile := strings.NewReader(`---
- version: 1
  output: "somenewfile.yml"
  imports:
  - "` + controlPathBase + `/local/file"`)

	noBaseYassemblyfile := strings.NewReader(`---
- version: 1
  output: "somenewfile.yml"
  imports:
  - "` + controlPathBase + `/local/file"`)

	t.Run("parse success", func(t *testing.T) {
		for _, table := range []struct {
			testName     string
			fileReader   io.Reader
			checkPatches bool
		}{
			{"defined single local import w/ patch", singlePatch, true},
			{"defined single local import", singleImport, false},
			{"defined multiple local import", multipleImports, false},
		} {
			t.Run(table.testName, func(t *testing.T) {
				configs, err := yassembly.ConfigParse(table.fileReader)
				if len(configs) == 0 {
					t.Error("Expected to have configs")
				}
				for _, config := range configs {
					for _, importPath := range config.Imports {
						if !strings.HasPrefix(importPath, controlPathBase) {
							t.Errorf("expected path prefix to be %s in path %s", controlPathBase, importPath)
						}
						if table.checkPatches && len(config.Patches) <= 0 {
							t.Error("expected to have patches")
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
			{"invalid y-assembly file", invalidYassemblyfile},
			{"no imports", noImportsYassemblyfile},
			{"no output yaml", noOutputYassemblyfile},
			{"no base yaml", noBaseYassemblyfile},
			{"empty y-assembly file", emptyYassemblyfile},
		} {
			t.Run(table.testName, func(t *testing.T) {
				configs, err := yassembly.ConfigParse(table.fileReader)
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
