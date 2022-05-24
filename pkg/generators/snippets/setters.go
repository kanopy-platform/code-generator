package snippets

import (
	"fmt"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

type Setter struct {
	Root            *types.Type
	Parent          *types.Type
	pointerReceiver bool
}

func NewSetter(root, parent *types.Type, pointerReceiver bool) *Setter {
	return &Setter{
		Root:            root,
		Parent:          parent,
		pointerReceiver: pointerReceiver,
	}
}

func (s *Setter) GenerateSetterForType(member types.Member) (string, generator.Args) {
	args := defaultGeneratorArgs(s.Root, s.pointerReceiver)
	args["funcName"] = funcName(member)
	args["memberAccessor"] = s.memberAccessor(member)
	args["memberType"] = member.Type

	raw := `// $.funcName$ is an autogenerated function
func (o $.pointer$$.type|raw$) $.funcName$(in $.memberType|raw$) $.pointer$$.type|raw$ {
	o.$.memberAccessor$ = in
	return o
}

`
	return raw, args
}

func (s *Setter) GenerateSetterForBool(member types.Member) (string, generator.Args) {
	args := defaultGeneratorArgs(s.Root, s.pointerReceiver)
	args["funcName"] = funcName(member)
	args["memberAccessor"] = s.memberAccessor(member)
	args["memberType"] = member.Type

	raw := `// $.funcName$ is an autogenerated function
func (o $.pointer$$.type|raw$) $.funcName$(in ...$.memberType|raw$) $.pointer$$.type|raw$ {
	o.$.memberAccessor$ = variadicBool(in...)
	return o
}

`
	return raw, args
}

func (s *Setter) GenerateSetterForPointerToBool(member types.Member) (string, generator.Args) {
	args := defaultGeneratorArgs(s.Root, s.pointerReceiver)
	args["funcName"] = funcName(member)
	args["memberAccessor"] = s.memberAccessor(member)
	args["memberType"] = member.Type.Elem

	raw := `// $.funcName$ is an autogenerated function
func (o $.pointer$$.type|raw$) $.funcName$(in ...$.memberType|raw$) $.pointer$$.type|raw$ {
	o.$.memberAccessor$ = boolPointer(variadicBool(in...))
	return o
}

`
	return raw, args
}

func (s *Setter) GenerateSetterForMap(member types.Member) (string, generator.Args) {
	args := defaultGeneratorArgs(s.Root, s.pointerReceiver)
	args["funcName"] = funcName(member)
	args["memberAccessor"] = s.memberAccessor(member)
	args["memberType"] = member.Type

	raw := `// $.funcName$ is an autogenerated function
func (o $.pointer$$.type|raw$) $.funcName$(in $.memberType|raw$) $.pointer$$.type|raw$ {
	if o.$.memberAccessor$ == nil {
		o.$.memberAccessor$ = make($.memberType|raw$)
	}
	for key, value := range in {
		o.$.memberAccessor$[key] = value
	}
	return o
}

`
	return raw, args
}

func (s *Setter) GenerateSetterForMapStringString(member types.Member) (string, generator.Args) {
	args := defaultGeneratorArgs(s.Root, s.pointerReceiver)
	args["funcName"] = funcName(member)
	args["memberAccessor"] = s.memberAccessor(member)
	args["memberType"] = member.Type

	raw := `// $.funcName$ is an autogenerated function
func (o $.pointer$$.type|raw$) $.funcName$(in $.memberType|raw$) $.pointer$$.type|raw$ {
	o.$.memberAccessor$ = mergeMapStringString(o.$.memberAccessor$, in)
	return o
}

`
	return raw, args
}

func (s *Setter) GenerateSetterForMemberSlice(member types.Member) (string, generator.Args) {
	args := defaultGeneratorArgs(s.Root, s.pointerReceiver)
	args["funcName"] = funcName(member)
	args["memberAccessor"] = s.memberAccessor(member)

	var raw string

	switch member.Type.Elem {
	case types.Byte:
		args["memberType"] = member.Type
		raw = `// $.funcName$ is an autogenerated function
func (o $.pointer$$.type|raw$) $.funcName$(in $.memberType|raw$) $.pointer$$.type|raw$ {
	o.$.memberAccessor$ = in
	return o
}

`
	default:
		args["memberType"] = member.Type.Elem
		raw = `// $.funcName$ is an autogenerated function
func (o $.pointer$$.type|raw$) $.funcName$(in ...$.memberType|raw$) $.pointer$$.type|raw$ {
	o.$.memberAccessor$ = append(o.$.memberAccessor$, in...)
	return o
}

`
	}

	return raw, args
}

func (s *Setter) GenerateSetterForEmbeddedSlice(member types.Member, inputType *types.Type) (string, generator.Args) {
	args := defaultGeneratorArgs(s.Root, s.pointerReceiver)
	args["funcName"] = funcName(member)
	args["memberAccessor"] = s.memberAccessor(member)
	args["inputType"] = inputType
	args["sliceType"] = member.Type.Elem.Name.Name

	raw := `// $.funcName$ is an autogenerated function
func (o $.pointer$$.type|raw$) $.funcName$(in ...*$.inputType|raw$) $.pointer$$.type|raw$ {
	for _, elem := range in {
		if elem != nil {
			o.$.memberAccessor$ = append(o.$.memberAccessor$, elem.$.sliceType$)
		}
	}
	return o
}

`
	return raw, args
}

func (s *Setter) GenerateSetterForMemberStruct(member types.Member) (string, generator.Args) {
	args := defaultGeneratorArgs(s.Root, s.pointerReceiver)
	args["funcName"] = funcName(member)
	args["memberAccessor"] = s.memberAccessor(member)
	args["memberType"] = member.Type

	raw := `// $.funcName$ is an autogenerated function
func (o $.pointer$$.type|raw$) $.funcName$(in *$.memberType|raw$) $.pointer$$.type|raw$ {
	if in != nil {
		o.$.memberAccessor$ = *in
	}
	return o
}

`
	return raw, args
}

func (s *Setter) GenerateSetterForEmbeddedStruct(member types.Member, inputType *types.Type) (string, generator.Args) {
	args := defaultGeneratorArgs(s.Root, s.pointerReceiver)
	args["funcName"] = funcName(member)
	args["memberAccessor"] = s.memberAccessor(member)
	args["inputType"] = inputType
	args["structType"] = member.Type.Name.Name

	raw := `// $.funcName$ is an autogenerated function
func (o $.pointer$$.type|raw$) $.funcName$(in *$.inputType|raw$) $.pointer$$.type|raw$ {
	if in != nil {
		o.$.memberAccessor$ = in.$.structType$
	}
	return o
}

`
	return raw, args
}

func (s *Setter) GenerateSetterForPointerToBuiltinType(member types.Member) (string, generator.Args) {
	args := defaultGeneratorArgs(s.Root, s.pointerReceiver)
	args["funcName"] = funcName(member)
	args["memberAccessor"] = s.memberAccessor(member)
	args["memberElemType"] = member.Type.Elem

	raw := `// $.funcName$ is an autogenerated function
func (o $.pointer$$.type|raw$) $.funcName$(in $.memberElemType|raw$) $.pointer$$.type|raw$ {
	o.$.memberAccessor$ = &in
	return o
}

`

	return raw, args
}

func (s *Setter) GenerateSetterForEmbeddedPointer(member types.Member, inputType *types.Type) (string, generator.Args) {
	args := defaultGeneratorArgs(s.Root, s.pointerReceiver)
	args["funcName"] = funcName(member)
	args["memberAccessor"] = s.memberAccessor(member)
	args["inputType"] = inputType
	args["structType"] = member.Type.Elem.Name.Name

	raw := `// $.funcName$ is an autogenerated function
func (o $.pointer$$.type|raw$) $.funcName$(in *$.inputType|raw$) $.pointer$$.type|raw$ {
	if in != nil {
		o.$.memberAccessor$ = &in.$.structType$
	}
	return o
}

`
	return raw, args
}

func funcName(m types.Member) string {
	verb := "With"

	if m.Type.Kind == types.Slice && m.Type.Elem != types.Byte {
		verb = "Append"
	}

	return fmt.Sprintf("%s%s", verb, m.Name)
}

func (s *Setter) memberAccessor(member types.Member) string {
	if s.Root != s.Parent {
		return fmt.Sprintf("%s.%s", s.Parent.Name.Name, member.Name)
	}
	return member.Name
}
