package tags

import (
	"strings"

	"k8s.io/gengo/types"
)

const (
	Builder        = "kanopy:builder"
	BuilderPackage = "package"
	BuilderOptIn   = "true"
	BuilderOptOut  = "false"
)

func IsPackageTagged(comments []string) bool {
	return Extract(comments, Builder) == BuilderPackage
}

func IsTypeEnabled(t *types.Type) bool {
	return Extract(combineTypeComments(t), Builder) == BuilderOptIn
}

func IsTypeOptedOut(t *types.Type) bool {
	return Extract(combineTypeComments(t), Builder) == BuilderOptOut
}

func IsMemberReadyOnly(m types.Member) bool {
	for _, s := range m.CommentLines {
		if strings.Contains(s, "Read-only") {
			return true
		}
	}
	return false
}

func Extract(comments []string, tag string) string {
	vals := types.ExtractCommentTags("+", comments)[tag]
	if len(vals) == 0 {
		return ""
	}

	return getFirstTagValue(vals...)
}

func combineTypeComments(t *types.Type) []string {
	return append(append([]string{}, t.SecondClosestCommentLines...), t.CommentLines...)
}

func getFirstTagValue(values ...string) string {
	return strings.Split(values[0], ",")[0]
}
