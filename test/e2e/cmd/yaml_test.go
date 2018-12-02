package main_test

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestYaml(t *testing.T) {
	gomega.RegisterTestingT(t)
	pathToYamlCLI, err := gexec.Build("github.com/xchapter7x/y-assembly/cmd/yaml")
	defer gexec.CleanupBuildArtifacts()
	defer func() {
		os.RemoveAll("testdata/outputs")
	}()
	if err != nil {
		t.Fatalf("build failed: %v", err)
	}

	t.Run("yaml build -c ./testdata/assembly.yml", func(t *testing.T) {
		command := exec.Command(pathToYamlCLI,
			"build",
			"-c",
			"./testdata/assembly.yml")
		session, err := gexec.Start(command, os.Stdout, os.Stderr)
		if err != nil {
			t.Fatalf("failed running command: %v", err)
		}
		session.Wait(120 * time.Second)
		if session.ExitCode() != 0 {
			t.Errorf("call failed: %v %v %v",
				session.ExitCode(),
				string(session.Out.Contents()),
				string(session.Err.Contents()))
		}
		importsControlBuffer := new(bytes.Buffer)
		baseFile, err := os.Open("testdata/base/base.yml")
		defer baseFile.Close()
		if err != nil {
			t.Fatal(err)
		}

		if _, err := io.Copy(importsControlBuffer, baseFile); err != nil {
			t.Fatal(err)
		}

		if _, err := io.Copy(importsControlBuffer, strings.NewReader("\n")); err != nil {
			t.Fatal(err)
		}

		importFile, err := os.Open("testdata/imports/import1.yml")
		defer importFile.Close()
		if err != nil {
			t.Fatal(err)
		}

		if _, err := io.Copy(importsControlBuffer, importFile); err != nil {
			t.Fatal(err)
		}
		patchesControlBuffer := bytes.NewBufferString(`lots:
- of
- stuff
- with
- even
- more
- imports
somecool: thing
which: had
`)

		for _, table := range []struct {
			name       string
			control    *bytes.Buffer
			outputPath string
		}{
			{"imports from remote source", importsControlBuffer, "testdata/outputs/out1.yml"},
			{"imports from local source", importsControlBuffer, "testdata/outputs/out2.yml"},
			{"imports with patch", patchesControlBuffer, "testdata/outputs/out3.yml"},
		} {

			t.Run(table.name, func(t *testing.T) {

				testBuffer := new(bytes.Buffer)
				outputFile, err := os.Open(table.outputPath)
				defer outputFile.Close()
				if err != nil {
					t.Fatal(err)
				}
				if _, err := io.Copy(testBuffer, outputFile); err != nil {
					t.Fatal(err)
				}
				if table.control.String() != testBuffer.String() {
					t.Errorf("generated file does not match the control: \n'%s' \n!=\n \n'%s'",
						table.control.String(),
						testBuffer.String(),
					)
				}
			})
		}
	})
}
