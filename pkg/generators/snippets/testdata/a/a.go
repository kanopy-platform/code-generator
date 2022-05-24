package a

import (
	"github.com/kanopy-platform/code-generator/pkg/generators/snippets/testdata/b"
)

type SomeStruct struct {
	b.TypeMeta
	b.ObjectMeta
	AStruct            AStruct
	CStruct            CStruct
	PointerCStruct     *CStruct
	CStructs           []CStruct
	Strings            []string
	Bytes              []byte
	IntPtr             *int
	MapIntString       map[int]string
	MapStringByteSlice map[string][]byte
	Bool               bool
	PointerBool        *bool
}

type AStruct struct {
	Value int
}

type CStruct struct {
	Int int
}
