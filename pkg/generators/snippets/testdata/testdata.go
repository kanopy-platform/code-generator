package testdata

import (
	"github.com/kanopy-platform/code-generator/pkg/generators/snippets/testdata/a"
)

type SomeStruct struct {
	a.SomeStruct
}

type CStruct struct {
	a.CStruct
}

type CopyStruct struct {
	a.SomeStruct
	ComplexStruct
}

type ComplexStruct struct {
	Strings []string
	IntPtr  *int
}

type Alias b.AliasOfString

func (c *ComplexStruct) DeepCopyInto(in *ComplexStruct) {
	// not impl test only
}
