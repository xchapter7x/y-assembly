package main_test

import (
	"crypto/md5"
	"fmt"
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
		hashControl := md5.New()
		baseFile, err := os.Open("testdata/base/base.yml")
		defer baseFile.Close()
		if err != nil {
			t.Fatal(err)
		}

		if _, err := io.Copy(hashControl, baseFile); err != nil {
			t.Fatal(err)
		}

		if _, err := io.Copy(hashControl, strings.NewReader("\n")); err != nil {
			t.Fatal(err)
		}

		importFile, err := os.Open("testdata/imports/import1.yml")
		defer importFile.Close()
		if err != nil {
			t.Fatal(err)
		}

		if _, err := io.Copy(hashControl, importFile); err != nil {
			t.Fatal(err)
		}

		hashGenerated := md5.New()
		outputFile, err := os.Open("testdata/outputs/out1.yml")
		defer outputFile.Close()
		if err != nil {
			t.Fatal(err)
		}
		if _, err := io.Copy(hashGenerated, outputFile); err != nil {
			t.Fatal(err)
		}
		hashControlString := fmt.Sprintf("%x", hashControl.Sum(nil))
		hashGeneratedString := fmt.Sprintf("%x", hashGenerated.Sum(nil))
		if hashControlString != hashGeneratedString {
			t.Errorf("generated file does not match the fixture hash: %s != %s",
				hashControlString,
				hashGeneratedString,
			)
		}
	})
}
