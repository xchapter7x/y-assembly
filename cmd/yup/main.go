package main

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/xchapter7x/yup/pkg/yup"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	Version   = "dev-build"
	Buildtime = "sometime"
	Platform  = "somecomputer"
	yupConfig = kingpin.Flag("config", "config file path").Short('c').Default(".yupfile").String()
	build     = kingpin.Command("build", "Build a yaml with your defined imports")
	version   = kingpin.Command("version", "display version information")
)

func main() {
	switch kingpin.Parse() {
	case version.FullCommand():
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Built On: %s\n", Buildtime)
		fmt.Printf("Platform: %s\n", Platform)
	case build.FullCommand():
		yupfile, err := os.Open(*yupConfig)
		if err != nil {
			panic(err)
		}
		configs, err := yup.ConfigParse(yupfile)
		if err != nil {
			panic(err)
		}
		for _, config := range configs {
			err := os.MkdirAll(path.Dir(config.Output), os.ModePerm)
			if err != nil {
				panic(err)
			}

			output, err := os.Create(config.Output)
			if err != nil {
				panic(err)
			}

			baseFile, err := os.Open(config.Base)
			if err != nil {
				panic(err)
			}

			imports := make([]io.Reader, 0)

			for _, importPath := range config.Imports {
				reader, err := os.Open(importPath)
				if err != nil {
					panic(err)
				}
				imports = append(imports, reader)
			}

			err = yup.Merge(baseFile, imports, output)
			if err != nil {
				panic(err)
			}
		}
	}
}
