package d

import (
	"github.com/kanopy-platform/code-generator/pkg/generators/builder/testdata/c/meta"
)

type MockDeployment struct {
	meta.TypeMeta
	meta.ObjectMeta
	Spec               MockSpec
	PointerSpec        *MockSpec
	Specs              []MockSpec
	SpecNoGen          MockSpecNoGen
	PointerSpecNoGen   *MockSpecNoGen
	Primitive          int
	Bool               bool
	PointerBool        *bool
	MapStringByteSlice map[string][]byte
}

type MockSpec struct {
}

type MockSpecNoGen struct {
}
