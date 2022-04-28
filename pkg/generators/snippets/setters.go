package snippets

import (
	"fmt"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

type Setter struct {
	Root   *types.Type
	Parent *types.Type
}

func NewSetter(root, parent *types.Type) *Setter {
	return &Setter{
		Root:   root,
		Parent: parent,
	}
}

func (s *Setter) GenerateSetterForPrimitiveType(member types.Member) (string, generator.Args) {
	args := generator.Args{
		"type":           s.Root,
		"funcName":       funcName(member),
		"memberAccessor": s.memberAccessor(member),
		"memberType":     member.Type,
	}

	raw := `// $.funcName$ is an autogenerated function
func (o *$.type|raw$) $.funcName$(in $.memberType|raw$) *$.type|raw$ {
	o.$.memberAccessor$ = in
	return o
}

`
	return raw, args
}

func (s *Setter) GenerateSetterForMap(member types.Member) (string, generator.Args) {
	args := generator.Args{
		"type":           s.Root,
		"funcName":       funcName(member),
		"memberAccessor": s.memberAccessor(member),
		"memberType":     member.Type,
	}

	raw := `// $.funcName$ is an autogenerated function
func (o *$.type|raw$) $.funcName$(in $.memberType|raw$) *$.type|raw$ {
	o.$.memberAccessor$ = mergeMapStringString(o.$.memberAccessor$, in)
	return o
}

`
	return raw, args
}

func (s *Setter) GenerateSetterForMemberSlice(member types.Member) (string, generator.Args) {
	args := generator.Args{
		"type":           s.Root,
		"funcName":       funcName(member),
		"memberAccessor": s.memberAccessor(member),
		"memberType":     member.Type.Elem,
	}

	raw := `// $.funcName$ is an autogenerated function
func (o *$.type|raw$) $.funcName$(in ...$.memberType|raw$) *$.type|raw$ {
	o.$.memberAccessor$ = append(o.$.memberAccessor$, in...)
	return o
}

`
	return raw, args
}

func (s *Setter) GenerateSetterForEmbeddedSlice(member types.Member, inputType *types.Type) (string, generator.Args) {
	args := generator.Args{
		"type":           s.Root,
		"funcName":       funcName(member),
		"memberAccessor": s.memberAccessor(member),
		"inputType":      inputType,
		"sliceType":      member.Type.Elem.Name.Name,
	}

	raw := `// $.funcName$ is an autogenerated function
func (o *$.type|raw$) $.funcName$(in ...$.inputType|raw$) *$.type|raw$ {
	for _, elem := range in {
		o.$.memberAccessor$ = append(o.$.memberAccessor$, elem.$.sliceType$)
	}
	return o
}

`
	return raw, args
}

func (s *Setter) GenerateSetterForMemberStruct(member types.Member) (string, generator.Args) {
	args := generator.Args{
		"type":           s.Root,
		"funcName":       funcName(member),
		"memberAccessor": s.memberAccessor(member),
		"memberType":     member.Type,
	}

	raw := `// $.funcName$ is an autogenerated function
func (o *$.type|raw$) $.funcName$(in $.memberType|raw$) *$.type|raw$ {
	o.$.memberAccessor$ = in
	return o
}

`
	return raw, args
}

func (s *Setter) GenerateSetterForEmbeddedStruct(member types.Member, inputType *types.Type) (string, generator.Args) {
	args := generator.Args{
		"type":           s.Root,
		"funcName":       funcName(member),
		"memberAccessor": s.memberAccessor(member),
		"inputType":      inputType,
		"structType":     member.Type.Name.Name,
	}

	raw := `// $.funcName$ is an autogenerated function
func (o *$.type|raw$) $.funcName$(in $.inputType|raw$) *$.type|raw$ {
	o.$.memberAccessor$ = in.$.structType$
	return o
}

`
	return raw, args
}

func (s *Setter) GenerateSetterForPointerToBuiltinType(member types.Member) (string, generator.Args) {
	args := generator.Args{
		"type":           s.Root,
		"funcName":       funcName(member),
		"memberAccessor": s.memberAccessor(member),
		"memberElemType": member.Type.Elem,
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

func (s *Setter) memberAccessor(member types.Member) string {
	if s.Root != s.Parent {
		return fmt.Sprintf("%s.%s", s.Parent.Name.Name, member.Name)
	}
	return member.Name
}
