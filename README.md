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
