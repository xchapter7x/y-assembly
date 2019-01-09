# y-assembly
YAML Plus Imports: y-assembly

[![CircleCI](https://circleci.com/gh/xchapter7x/y-assembly.svg?style=svg)](https://circleci.com/gh/xchapter7x/y-assembly)

### Description

This tool is meant to provide some generic "imports" or merging capability.
One can combine files and apply patch files which are yaml patch compliant.
Combining this functionality with yaml references and/or envsubst, makes for a really powerful and 
flexible combination.

### Use case

anywhere you are wrangling yaml (Kubernetes, BOSH, Concourse) this can help. It is generic, specifically built to provide enough flexibility to be used with any tool that likes yaml.


### Assembly.yml fields and file behavior

`version`: and integer value representing the compatibility of the assembly record 

`base`: the base file, where all others are appended or modify (can be url or local file)

`output`: where the output of the assembly operation will go (must be local)

`imports`: an array of urls or local file paths, which will be appended to the base to create the output. the array order is the order in which the files will be combined.

`patches`: an array of urls or local file paths, which will be patched on top of the fully assembled set of base and imports. the array order is the order in which the patches will be applied. uses [yaml-patch op files](https://github.com/krishicks/yaml-patch) for full set of commands check the compatible [json patch docs](https://tools.ietf.org/html/rfc6902)


### Example
```bash
-> % cd github.com/xchapter7x/y-assembly/test/e2e/cmd

-> % cat testdata/assembly.yml
---
- version: 1
  base: "testdata/base/base.yml"
  output: "testdata/outputs/cool.yml"
  imports:
  - "https://raw.githubusercontent.com/xchapter7x/y-assembly/master/test/e2e/cmd/testdata/imports/import1.yml"

- version: 1
  base: "testdata/base/base.yml"
  output: "testdata/outputs/out2.yml"
  imports:
  - "testdata/imports/import1.yml"

- version: 1
  base: "testdata/base/base.yml"
  output: "testdata/outputs/out3.yml"
  imports:
  - "https://raw.githubusercontent.com/xchapter7x/y-assembly/master/test/e2e/cmd/testdata/imports/import1.yml"
  patches:
  - "testdata/patches/patch1.yml"
  
-> % ./yaml_osx build -c testdata/assembly.yml -p

---
#  testdata/outputs/out1.yml
somecool: thing
which: has
lots:
- of
- stuff

- with
- even
- more
- imports

---
#  testdata/outputs/out2.yml
somecool: thing
which: has
lots:
- of
- stuff

- with
- even
- more
- imports

---
#  testdata/outputs/out3.yml
lots:
- of
- stuff
- with
- even
- more
- imports
somecool: thing
which: had
```
