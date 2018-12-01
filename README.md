# y-assembly
YAML Plus Imports: y-assembly

### Description

This tool is meant to provide some generic "imports" or merging capability. Inent is to 
use reference links to yaml blocks which reside in an "import" file. This way
one can inject different implementations via different imports.

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

`imports`: an array of urls or local file paths, which will be appended to the base to create the output 
