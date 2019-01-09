package yassembly_test

import (
	"bytes"
	"testing"

	yassembly "github.com/xchapter7x/y-assembly/pkg/y-assembly"
)

func TestExpandAliases(t *testing.T) {
	controlAliases := bytes.NewReader([]byte(`---
mylist: &mylist
  some: thingelse
  other: randomthing 
name: somthing
age:
  << : *mylist
`))
	testOutput := new(bytes.Buffer)

	controlOutput := []byte(`age:
  other: randomthing
  some: thingelse
mylist:
  other: randomthing
  some: thingelse
name: somthing
`)

	err := yassembly.ExpandAliases(controlAliases, testOutput)
	if err != nil {
		t.Errorf("should not error on valid inputs: %v", err)
	}

	if string(testOutput.Bytes()) != string(controlOutput) {
		t.Errorf(
			"output should match: \n'%v'\n!=\n'%v'\n",
			string(testOutput.Bytes()),
			string(controlOutput),
		)
	}
}
