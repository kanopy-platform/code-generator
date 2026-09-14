package e

import (
	"github.com/kanopy-platform/code-generator/pkg/generators/builder/testdata/c/d"
)

// +kanopy:builder=true
type EDeployment struct {
	d.MockDeployment
}
