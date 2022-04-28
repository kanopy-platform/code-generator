package d

import (
	"github.com/kanopy-platform/code-generator/pkg/generators/builder/testdata/c/meta"
)

type MockDeployment struct {
	meta.TypeMeta
	meta.ObjectMeta
	Spec      MockSpec
	Specs     []MockSpec
	SpecNoGen MockSpecNoGen
	Primitive bool
}

type MockSpec struct {
}

type MockSpecNoGen struct {
}
