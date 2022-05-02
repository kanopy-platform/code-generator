package c

import (
	"github.com/kanopy-platform/code-generator/pkg/generators/builder/testdata/c/d"
)

// +kanopy:builder=true
type CDeployment struct {
	d.MockDeployment
}

// +kanopy:builder=true
// +kanopy:receiver=value
type MockSpec struct {
	d.MockSpec
}
