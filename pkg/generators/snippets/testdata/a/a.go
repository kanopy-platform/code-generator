package a

import (
	"github.com/kanopy-platform/code-generator/pkg/generators/snippets/testdata/b"
)

type SomeStruct struct {
	b.TypeMeta
	b.ObjectMeta
	AStruct            AStruct
	CStruct            CStruct
	CStructs           []CStruct
	Strings            []string
	IntPtr             *int
	MapIntString       map[int]string
	MapStringByteSlice map[string][]byte
}

type AStruct struct {
	Value int
}

type CStruct struct {
	Int int
}
