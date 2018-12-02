# y-assembly
YAML Plus Imports: y-assembly

### Description

This tool is meant to provide some generic "imports" or merging capability. Inent is to 
use reference links to yaml blocks which reside in an "import" file. This way
one can inject different implementations via different imports.

### Use case

anywhere you are wrangling yaml (Kubernetes, BOSH, Concourse) this can help. It is generic, specifically built to provide enough flexibility to be used with any tool that likes yaml.

### Example
```yaml
# basefile.yml
---
mycool: file
somelist:
<<: *mylist

# importlist_A.yml
mylist: &mylist
- some
- set
- of 
- data

# importlist_B.yml
mylist: &mylist
- some
- other 
- random
- dataset

# assembly.yml
- version: 1
  base: basefile.yml
  output: cool.yml
  imports: 
  - importlist_B.yml
```


```bash
$ yaml build
$ cat cool.yml

---
mycool: file
somelist:
<<: *mylist

mylist: &mylist
- some
- other 
- random
- dataset
```

### Assembly.yml fields and file behavior

`version`: and integer value representing the compatibility of the assembly record 

`base`: the base file, where all others are appended or modify (can be url or local file)

`output`: where the output of the assembly operation will go (must be local)

`imports`: an array of urls or local file paths, which will be appended to the base to create the output. the array order is the order in which the files will be combined.

`patches`: an array of urls or local file paths, which will be patched on top of the fully assembled set of base and imports. the array order is the order in which the patches will be applied. uses [yaml-patch op files](https://github.com/krishicks/yaml-patch) for full set of commands check the compatible [json patch docs](https://tools.ietf.org/html/rfc6902)
