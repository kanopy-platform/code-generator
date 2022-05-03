package a

import (
	"github.com/kanopy-platform/code-generator/pkg/generators/snippets/testdata/b"
)

type SomeStruct struct {
	b.TypeMeta
	b.ObjectMeta
	AStruct  AStruct
	CStruct  CStruct
	CStructs []CStruct
	Strings  []string
	IntPtr   *int
}

type AStruct struct {
	Value int
}

type CStruct struct {
	Int int
}
