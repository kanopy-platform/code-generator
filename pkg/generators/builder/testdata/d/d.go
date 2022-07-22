package d

import (
	"github.com/kanopy-platform/code-generator/pkg/generators/builder/testdata/d/e"
)

// +kanopy:builder=true
type DPolicyRule struct {
	e.MockPolicyRule
}

// +kanopy:builder=true
type AliasType e.AliasToString
