package d

import (
	"github.com/kanopy-platform/code-generator/pkg/generators/builder/testdata/d/e"
)

// +kanopy:builder=true
// +kanopy:receiver=value
type DPolicyRule struct {
	e.MockPolicyRule
}
