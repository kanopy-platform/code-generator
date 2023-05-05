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
	EnumFlag       = "enum"
	RefFlag        = "ref"
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

func GetEnumOptions(t *types.Type) []string {
	val := ExtractArg(combineTypeComments(t), Builder, EnumFlag)
	return strings.Split(val, ";")
}

func ExtractRef(t *types.Type) string {
	return ExtractArg(combineTypeComments(t), Builder, RefFlag)
}

func IsMemberReadyOnly(m types.Member) bool {
	for _, s := range m.CommentLines {
		if strings.Contains(strings.ToLower(s), "read-only") {
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

func ExtractArg(comments []string, tag string, arg string) string {
	vals := types.ExtractCommentTags("+", comments)[tag]
	if len(vals) == 0 {
		return ""
	}

	args := strings.Split(vals[0], ",")
	for _, a := range args {
		if strings.Contains(a, "=") {
			kvs := strings.Split(a, "=")
			if len(kvs) == 2 {
				key := kvs[0]
				if key == arg {
					return kvs[1]
				}
			}
		}

		if a == arg {
			return a
		}
	}
	return ""
}

func combineTypeComments(t *types.Type) []string {
	return append(append([]string{}, t.SecondClosestCommentLines...), t.CommentLines...)
}

func getFirstTagValue(values ...string) string {
	return strings.Split(values[0], ",")[0]
}
