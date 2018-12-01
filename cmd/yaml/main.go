package main

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/xchapter7x/y-assembly/pkg/y-assembly"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	Version         = "dev-build"
	Buildtime       = "sometime"
	Platform        = "somecomputer"
	yassemblyConfig = kingpin.Flag("config", "config file path").Short('c').Default("assembly.yml").String()
	build           = kingpin.Command("build", "Build a yaml with your defined imports")
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

			output, err := os.Create(config.Output)
			defer output.Close()
			if err != nil {
				panic(err)
			}

			baseFile, err := os.Open(config.Base)
			defer baseFile.Close()
			if err != nil {
				panic(err)
			}

			imports := make([]io.Reader, 0)

			for _, importPath := range config.Imports {
				importFile, err := os.Open(importPath)
				defer importFile.Close()
				if err != nil {
					panic(err)
				}
				imports = append(imports, importFile)
			}

			err = yassembly.Merge(baseFile, imports, output)
			if err != nil {
				panic(err)
			}
		}
	}
}
