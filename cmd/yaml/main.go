package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"

	"github.com/xchapter7x/y-assembly/pkg/y-assembly"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	Version         = "dev-build"
	Buildtime       = "sometime"
	Platform        = "somecomputer"
	yassemblyConfig = kingpin.Flag("config", "config file path").Short('c').Default("assembly.yml").String()
	build           = kingpin.Command("build", "Build a yaml with your defined imports")
	printYaml       = build.Flag("print", "Print yaml to STDOUT").Short('p').Bool()
	version         = kingpin.Command("version", "display version information")
)

func main() {
	switch kingpin.Parse() {
	case version.FullCommand():
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Built On: %s\n", Buildtime)
		fmt.Printf("Platform: %s\n", Platform)
	case build.FullCommand():
		yassemblyfile, err := os.Open(*yassemblyConfig)
		defer yassemblyfile.Close()
		if err != nil {
			panic(err)
		}
		configs, err := yassembly.ConfigParse(yassemblyfile)
		if err != nil {
			panic(err)
		}
		for _, config := range configs {
			err := os.MkdirAll(path.Dir(config.Output), os.ModePerm)
			if err != nil {
				panic(err)
			}

			outputBuffer := new(bytes.Buffer)
			baseFile, err := open(config.Base)
			defer baseFile.Close()
			if err != nil {
				panic(err)
			}

			imports := make([]io.Reader, 0)
			for _, importPath := range config.Imports {
				importFile, err := open(importPath)
				defer importFile.Close()
				if err != nil {
					panic(err)
				}
				imports = append(imports, importFile)
			}

			err = yassembly.Combine(baseFile, imports, outputBuffer)
			if err != nil {
				panic(err)
			}

			patches := make([]io.Reader, 0)
			for _, patchPath := range config.Patches {
				patchFile, err := open(patchPath)
				defer patchFile.Close()
				if err != nil {
					panic(err)
				}
				patches = append(patches, patchFile)
			}

			patchedOutputBuffer := new(bytes.Buffer)
			err = yassembly.Patch(outputBuffer, patches, patchedOutputBuffer)
			if err != nil {
				panic(err)
			}

			outputFile, err := getWriter(config.Output, *printYaml)
			defer outputFile.Close()
			if err != nil {
				panic(err)
			}

			_, err = io.Copy(outputFile, patchedOutputBuffer)
			if err != nil {
				panic(err)
			}
		}
	}
}

func getWriter(outputPath string, shouldPrint bool) (io.WriteCloser, error) {
	if shouldPrint {
		fmt.Fprintln(os.Stdout, "\n---")
		fmt.Fprintln(os.Stdout, "# ", outputPath)
		return os.Stdout, nil
	}
	return os.Create(outputPath)
}

func open(path string) (io.ReadCloser, error) {
	if isURL, _ := regexp.MatchString("^http.*://", path); isURL {
		return downloadFile(path)
	}
	return os.Open(path)
}

func downloadFile(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		defer resp.Body.Close()
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	return resp.Body, nil
}
