package testdata

import (
	"github.com/kanopy-platform/code-generator/pkg/generators/snippets/testdata/a"
)

// +kanopy:receiver=pointer
type SomeStruct struct {
	a.MockStruct
}

// +kanopy:receiver=value
type ValueStruct struct {
	a.MockStruct
}

type CStruct struct {
	a.CStruct
}
