package snippets

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

func GenerateSetter(w io.Writer, ctx *generator.Context, root *types.Type, member types.Member, hasParent bool, embedded *types.Type) error {
	if root == nil {
		return errors.New("nil pointer")
	}

	args := generator.Args{
		"type":       root,
		"funcName":   funcName(member),
		"memberName": member.Name,
		"memberType": member.Type,
	}

	if embedded != nil {
		args["embeddedType"] = embedded
	}

	assignment := "o.$.memberName$"
	if hasParent {
		assignment = "o.$.type|raw$.$.memberName$"
	}

	var sb strings.Builder
	sb.WriteString("// $.funcName$ is an autogenerated setter function.\n")

	switch member.Type.Kind {
	case types.Map:
		sb.WriteString("func (o *$.type|raw$) $.funcName$(in $.memberType|raw$) *$.type|raw$ {\n")
		sb.WriteString(fmt.Sprintf("%s = mergeMapStringString(%s, in)\n", assignment, assignment))
	case types.Slice:
		sb.WriteString("func (o *$.type|raw$) $.funcName$(in ...$.memberType|raw$) *$.type|raw$ {\n")
		if embedded != nil {
			sb.WriteString("for _, i := range in {\n")
			sb.WriteString(fmt.Sprintf("%s = append(%s, i.$.embeddedType|raw$)\n", assignment, assignment))
			sb.WriteString("}\n")
		} else {
			sb.WriteString(fmt.Sprintf("%s = append(%s, in...)\n", assignment, assignment))
		}
	case types.Struct:
		sb.WriteString("func (o *$.type|raw$) $.funcName$(in $.memberType|raw$) *$.type|raw$ {\n")
		if embedded != nil {
			sb.WriteString(fmt.Sprintf("%s = in.$.embeddedType|raw$\n", assignment))
		} else {
			sb.WriteString(fmt.Sprintf("%s = in\n", assignment))
		}
	// TODO pointers. Handle case of builtin types (e.g. *int32) to convert internally
	default:
		sb.WriteString("func (o *$.type|raw$) $.funcName$(in $.memberType|raw$) *$.type|raw$ {\n")
		sb.WriteString(fmt.Sprintf("%s = in\n", assignment))
	}

	sb.WriteString("return o\n")
	sb.WriteString("}\n\n")

	return writeSnippet(w, ctx, sb.String(), args)
}

func funcName(m types.Member) string {
	verb := "With"

	if m.Type.Kind == types.Slice {
		verb = "Append"
	}

	return fmt.Sprintf("%s%s", verb, m.Name)
}
