package tags

import (
	"strings"

	"k8s.io/gengo/types"
)

const (
	Builder         = "kanopy:builder"
	Builder_Package = "package"
	Builder_OptIn   = "true"
	Builder_OptOut  = "false"
	Receiver        = "kanopy:receiver"
	Receiver_Ptr    = "ptr"
	Receiver_Value  = "value"
)

func IsPackageTagged(comments []string) bool {
	return Extract(comments, Builder) == Builder_Package
}

func IsTypeEnabled(t *types.Type) bool {
	return Extract(combineTypeComments(t), Builder) == Builder_OptIn
}

func IsTypeOptedOut(t *types.Type) bool {
	return Extract(combineTypeComments(t), Builder) == Builder_OptOut
}

func IsPtrReceiver(t *types.Type) bool {
	val := Extract(combineTypeComments(t), Receiver)
	// default to Pointer Receiver if unspecified
	return (val == Receiver_Ptr) || (val == "")
}

func IsValueReceiver(t *types.Type) bool {
	return Extract(combineTypeComments(t), Receiver) == Receiver_Value
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
