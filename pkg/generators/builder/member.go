package builder

import (
	"github.com/kanopy-platform/code-generator/pkg/generators/tags"
	log "github.com/sirupsen/logrus"
	"k8s.io/gengo/types"
)

const (
	ObjectMeta = "ObjectMeta"
)

func includeMember(parent *types.Type, member types.Member) bool {

	log.Debugf("includeMember Check %v", member.Name)

	if tags.IsMemberReadyOnly(member) {
		log.Debugf("\t member %v is readonly", member.Name)
		return false
	}

	switch parent.Name.Name {
	case ObjectMeta:
		log.Debug("\t member has ObjectMeta")
		return includeObjectMetaMember(member)
	default:
		log.Debug("\t included")
		return true
	}
}

func includeObjectMetaMember(member types.Member) bool {
	return member.Name != "Finalizers"
}
