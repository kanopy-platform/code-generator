package tags

import (
	"strings"

	"k8s.io/gengo/types"
)

const (
	Name    = "kanopy:builder"
	Package = "package"
)

func IsPackageTagged(tag string) bool {
	return tag == Package
}

func IsTypeEnabled(t *types.Type) bool {
	tag := Extract(combineTypeComments(t))
	return tag == "true"
}

func IsTypeOptedOut(t *types.Type) bool {
	tag := Extract(combineTypeComments(t))
	return tag == "false"
}

func Extract(comments []string) string {
	vals := types.ExtractCommentTags("+", comments)[Name]
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
