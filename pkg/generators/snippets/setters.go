package snippets

import (
	"fmt"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

func GenerateSetterPrimitive(root *types.Type, member types.Member) (string, generator.Args) {
	args := generator.Args{
		"type":       root,
		"funcName":   funcName(member),
		"memberName": member.Name,
		"memberType": member.Type,
	}

	raw := `// $.funcName$ is an autogenerated function
func (o *$.type|raw$) $.funcName$(in $.memberType|raw$) *$.type|raw$ {
	o.$.memberName$ = in
	return o
}

`
	return raw, args
}

func GenerateSetterMap(root *types.Type, member types.Member) (string, generator.Args) {
	args := generator.Args{
		"type":       root,
		"funcName":   funcName(member),
		"memberName": member.Name,
		"memberType": member.Type,
	}

	raw := `// $.funcName$ is an autogenerated function
func (o *$.type|raw$) $.funcName$(in $.memberType|raw$) *$.type|raw$ {
	o.$.memberName$ = mergeMapStringString(o.$.memberName$, in)
	return o
}

`
	return raw, args
}

func GenerateSetterSlice(root *types.Type, member types.Member, wrapper *types.Type) (string, generator.Args) {
	args := generator.Args{
		"type":       root,
		"funcName":   funcName(member),
		"memberName": member.Name,
	}

	var raw string
	if wrapper != nil {
		// needs translation from struct to pull out the embedded type
		args["wrapperType"] = wrapper
		args["embeddedTypeName"] = member.Type.Elem.Name.Name

		raw = `// $.funcName$ is an autogenerated function
func (o *$.type|raw$) $.funcName$(in ...$.wrapperType|raw$) *$.type|raw$ {
	for _, elem := range in {
		o.$.memberName$ = append(o.$.memberName$, elem.$.embeddedTypeName$)
	}
	return o
}

`
	} else {
		args["memberType"] = member.Type.Elem
		raw = `// $.funcName$ is an autogenerated function
func (o *$.type|raw$) $.funcName$(in ...$.memberType|raw$) *$.type|raw$ {
	o.$.memberName$ = append(o.$.memberName$, in...)
	return o
}

`
	}

	return raw, args
}

func GenerateSetterStruct(root *types.Type, member types.Member, wrapper *types.Type) (string, generator.Args) {
	args := generator.Args{
		"type":       root,
		"funcName":   funcName(member),
		"memberName": member.Name,
	}

	var raw string
	if wrapper != nil {
		// needs translation from struct to pull out the embedded type
		args["wrapperType"] = wrapper
		args["embeddedTypeName"] = member.Type.Name.Name

		raw = `// $.funcName$ is an autogenerated function
func (o *$.type|raw$) $.funcName$(in $.wrapperType|raw$) *$.type|raw$ {
	o.$.memberName$ = in.$.embeddedTypeName$
	return o
}

`
	} else {
		args["memberType"] = member.Type
		raw = `// $.funcName$ is an autogenerated function
func (o *$.type|raw$) $.funcName$(in $.memberType|raw$) *$.type|raw$ {
	o.$.memberName$ = in
	return o
}

`
	}

	return raw, args
}

func GenerateSetterPointerBuiltin(root *types.Type, member types.Member) (string, generator.Args) {
	args := generator.Args{
		"type":           root,
		"funcName":       funcName(member),
		"memberName":     member.Name,
		"memberElemType": member.Type.Elem,
	}

	raw := `// $.funcName$ is an autogenerated function
func (o *$.type|raw$) $.funcName$(in $.memberElemType|raw$) *$.type|raw$ {
	o.$.memberName$ = &in
	return o
}

`

	return raw, args
}

func funcName(m types.Member) string {
	verb := "With"

	if m.Type.Kind == types.Slice {
		verb = "Append"
	}

	return fmt.Sprintf("%s%s", verb, m.Name)
}
