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

	var err error

	if member.Type.Kind == types.Map {
		err = generateMapSetter(w, ctx, args, assignment)
	} else if member.Type.Kind == types.Slice {
		err = generateSliceSetter(w, ctx, args, assignment, embedded)
	} else if member.Type.Kind == types.Struct {
		err = generateStructSetter(w, ctx, args, assignment, embedded)
	} else if member.Type.Kind == types.Pointer && member.Type.Elem.Kind == types.Builtin {
		args["memberElemType"] = member.Type.Elem
		err = generatePointerBuiltin(w, ctx, args, assignment)
	} else {
		err = generateDefaultSetter(w, ctx, args, assignment)
	}

	return err
}

func funcName(m types.Member) string {
	verb := "With"

	if m.Type.Kind == types.Slice {
		verb = "Append"
	}

	return fmt.Sprintf("%s%s", verb, m.Name)
}

func commentHeader() string {
	return "// $.funcName$ is an autogenerated setter function.\n"
}

func returnFooter() string {
	return "return o\n}\n\n"
}

func generateDefaultSetter(w io.Writer, ctx *generator.Context, args generator.Args, assignment string) error {
	var sb strings.Builder
	sb.WriteString(commentHeader())

	sb.WriteString("func (o *$.type|raw$) $.funcName$(in $.memberType|raw$) *$.type|raw$ {\n")
	sb.WriteString(fmt.Sprintf("%s = in\n", assignment))

	sb.WriteString(returnFooter())
	return writeSnippet(w, ctx, sb.String(), args)
}

func generateMapSetter(w io.Writer, ctx *generator.Context, args generator.Args, assignment string) error {
	var sb strings.Builder
	sb.WriteString(commentHeader())

	sb.WriteString("func (o *$.type|raw$) $.funcName$(in $.memberType|raw$) *$.type|raw$ {\n")
	sb.WriteString(fmt.Sprintf("%s = mergeMapStringString(%s, in)\n", assignment, assignment))

	sb.WriteString(returnFooter())
	return writeSnippet(w, ctx, sb.String(), args)
}

func generateSliceSetter(w io.Writer, ctx *generator.Context, args generator.Args, assignment string, embedded *types.Type) error {
	var sb strings.Builder
	sb.WriteString(commentHeader())

	sb.WriteString("func (o *$.type|raw$) $.funcName$(in ...$.memberType|raw$) *$.type|raw$ {\n")
	if embedded != nil {
		// needs translation from struct to pull out the embeddedType
		sb.WriteString("for _, i := range in {\n")
		sb.WriteString(fmt.Sprintf("%s = append(%s, i.$.embeddedType|raw$)\n", assignment, assignment))
		sb.WriteString("}\n")
	} else {
		sb.WriteString(fmt.Sprintf("%s = append(%s, in...)\n", assignment, assignment))
	}

	sb.WriteString(returnFooter())
	return writeSnippet(w, ctx, sb.String(), args)
}

func generateStructSetter(w io.Writer, ctx *generator.Context, args generator.Args, assignment string, embedded *types.Type) error {
	var sb strings.Builder
	sb.WriteString(commentHeader())

	sb.WriteString("func (o *$.type|raw$) $.funcName$(in $.memberType|raw$) *$.type|raw$ {\n")
	if embedded != nil {
		// needs translation from struct to pull out the embeddedType
		sb.WriteString(fmt.Sprintf("%s = in.$.embeddedType|raw$\n", assignment))
	} else {
		sb.WriteString(fmt.Sprintf("%s = in\n", assignment))
	}

	sb.WriteString(returnFooter())
	return writeSnippet(w, ctx, sb.String(), args)
}

func generatePointerBuiltin(w io.Writer, ctx *generator.Context, args generator.Args, assignment string) error {
	var sb strings.Builder
	sb.WriteString(commentHeader())

	// convert from builtin type to pointer
	sb.WriteString("func (o *$.type|raw$) $.funcName$(in $.memberElemType|raw$) *$.type|raw$ {\n")
	sb.WriteString(fmt.Sprintf("%s = &in\n", assignment))

	sb.WriteString(returnFooter())
	return writeSnippet(w, ctx, sb.String(), args)
}
