package snippets

import (
	"fmt"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

func GenerateSetterForPrimitiveType(root *types.Type, parent *types.Type, member types.Member) (string, generator.Args) {
	args := generator.Args{
		"type":           root,
		"funcName":       funcName(member),
		"memberAccessor": member.Name,
		"memberType":     member.Type,
	}

	if root != parent {
		args["memberAccessor"] = fmt.Sprintf("%s.%s", parent.Name.Name, member.Name)
	}

	raw := `// $.funcName$ is an autogenerated function
func (o *$.type|raw$) $.funcName$(in $.memberType|raw$) *$.type|raw$ {
	o.$.memberAccessor$ = in
	return o
}

`
	return raw, args
}

func GenerateSetterForMap(root *types.Type, parent *types.Type, member types.Member) (string, generator.Args) {
	args := generator.Args{
		"type":           root,
		"funcName":       funcName(member),
		"memberAccessor": member.Name,
		"memberType":     member.Type,
	}

	if root != parent {
		args["memberAccessor"] = fmt.Sprintf("%s.%s", parent.Name.Name, member.Name)
	}

	raw := `// $.funcName$ is an autogenerated function
func (o *$.type|raw$) $.funcName$(in $.memberType|raw$) *$.type|raw$ {
	o.$.memberAccessor$ = mergeMapStringString(o.$.memberAccessor$, in)
	return o
}

`
	return raw, args
}

func GenerateSetterForMemberSlice(root *types.Type, parent *types.Type, member types.Member) (string, generator.Args) {
	args := generator.Args{
		"type":           root,
		"funcName":       funcName(member),
		"memberAccessor": member.Name,
		"memberType":     member.Type.Elem,
	}

	if root != parent {
		args["memberAccessor"] = fmt.Sprintf("%s.%s", parent.Name.Name, member.Name)
	}

	raw := `// $.funcName$ is an autogenerated function
func (o *$.type|raw$) $.funcName$(in ...$.memberType|raw$) *$.type|raw$ {
	o.$.memberAccessor$ = append(o.$.memberAccessor$, in...)
	return o
}

`
	return raw, args
}

func GenerateSetterForEmbeddedSlice(root *types.Type, parent *types.Type, member types.Member, wrapper *types.Type) (string, generator.Args) {
	args := generator.Args{
		"type":             root,
		"funcName":         funcName(member),
		"memberAccessor":   member.Name,
		"wrapper":          wrapper,
		"embeddedTypeName": member.Type.Elem.Name.Name,
	}

	if root != parent {
		args["memberAccessor"] = fmt.Sprintf("%s.%s", parent.Name.Name, member.Name)
	}

	raw := `// $.funcName$ is an autogenerated function
func (o *$.type|raw$) $.funcName$(in ...$.wrapper|raw$) *$.type|raw$ {
	for _, elem := range in {
		o.$.memberAccessor$ = append(o.$.memberAccessor$, elem.$.embeddedTypeName$)
	}
	return o
}

`
	return raw, args
}

func GenerateSetterForMemberStruct(root *types.Type, parent *types.Type, member types.Member) (string, generator.Args) {
	args := generator.Args{
		"type":           root,
		"funcName":       funcName(member),
		"memberAccessor": member.Name,
		"memberType":     member.Type,
	}

	if root != parent {
		args["memberAccessor"] = fmt.Sprintf("%s.%s", parent.Name.Name, member.Name)
	}

	raw := `// $.funcName$ is an autogenerated function
func (o *$.type|raw$) $.funcName$(in $.memberType|raw$) *$.type|raw$ {
	o.$.memberAccessor$ = in
	return o
}

`
	return raw, args
}

func GenerateSetterForEmbeddedStruct(root *types.Type, parent *types.Type, member types.Member, wrapper *types.Type) (string, generator.Args) {
	args := generator.Args{
		"type":             root,
		"funcName":         funcName(member),
		"memberAccessor":   member.Name,
		"wrapper":          wrapper,
		"embeddedTypeName": member.Type.Name.Name,
	}

	if root != parent {
		args["memberAccessor"] = fmt.Sprintf("%s.%s", parent.Name.Name, member.Name)
	}

	raw := `// $.funcName$ is an autogenerated function
func (o *$.type|raw$) $.funcName$(in $.wrapper|raw$) *$.type|raw$ {
	o.$.memberAccessor$ = in.$.embeddedTypeName$
	return o
}

`
	return raw, args
}

func GenerateSetterForPointerToBuiltinType(root *types.Type, parent *types.Type, member types.Member) (string, generator.Args) {
	args := generator.Args{
		"type":           root,
		"funcName":       funcName(member),
		"memberAccessor": member.Name,
		"memberElemType": member.Type.Elem,
	}

	if root != parent {
		args["memberAccessor"] = fmt.Sprintf("%s.%s", parent.Name.Name, member.Name)
	}

	raw := `// $.funcName$ is an autogenerated function
func (o *$.type|raw$) $.funcName$(in $.memberElemType|raw$) *$.type|raw$ {
	o.$.memberAccessor$ = &in
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
