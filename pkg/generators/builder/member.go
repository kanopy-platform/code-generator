package builder

import (
	"github.com/kanopy-platform/code-generator/pkg/generators/tags"
	"k8s.io/gengo/types"
)

func includeMember(parent *types.Type, member types.Member) bool {
	if tags.IsMemberReadyOnly(member) {
		return false
	}

	switch parent.Name.Name {
	case "ObjectMeta":
		return includeObjectMetaMember(member)
	default:
		return true
	}
}

func includeObjectMetaMember(member types.Member) bool {
	return member.Name != "Finalizers"
}
