package a

import (
	"github.com/kanopy-platform/code-generator/pkg/generators/snippets/testdata/b"
)

const (
	// Number is the index of the member in the struct
	MemberIndex_SomeStruct_TypeMeta   = 0
	MemberIndex_SomeStruct_ObjectMeta = 1
	MemberIndex_SomeStruct_AStruct    = 2
	MemberIndex_SomeStruct_CStruct    = 3
	MemberIndex_SomeStruct_CStructs   = 4
	MemberIndex_SomeStruct_Strings    = 5
	MemberIndex_SomeStruct_IntPtr     = 6
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
