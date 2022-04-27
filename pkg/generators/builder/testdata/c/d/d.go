package d

import (
	"github.com/kanopy-platform/code-generator/pkg/generators/builder/testdata/c/meta"
)

type MockDeployment struct {
	meta.TypeMeta
	meta.ObjectMeta
	Spec MockDeploymentSpec
}

type MockDeploymentSpec struct {
	Replicas int
}
