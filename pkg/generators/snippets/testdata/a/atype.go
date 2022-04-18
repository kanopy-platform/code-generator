package a

import (
	"github.com/kanopy-platform/code-generator/pkg/generators/snippets/testdata/b"
	"github.com/kanopy-platform/code-generator/pkg/generators/snippets/testdata/c"
)

const (
	// Match the number with the index of the member in the struct
	MemberIndex_AStruct_TypeMeta       = 0
	MemberIndex_AStruct_ObjectMeta     = 1
	MemberIndex_AStruct_ComplexStruct  = 2
	MemberIndex_AStruct_ComplexStructs = 3
	MemberIndex_AStruct_SimpleStruct   = 4
	MemberIndex_AStruct_Strings        = 5
	MemberIndex_AStruct_IntPtr         = 6
)

type AStruct struct {
	b.TypeMeta
	b.ObjectMeta
	ComplexStruct  ComplexStruct
	ComplexStructs []ComplexStruct
	SimpleStruct   SimpleStruct
	Strings        []string
	IntPtr         *int
}

type ComplexStruct struct {
	c.ComplexStruct
}

type SimpleStruct struct {
	Value int
}
