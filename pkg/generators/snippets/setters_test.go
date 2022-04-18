package snippets

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

func TestGenerateSetter(t *testing.T) {
	t.Parallel()

	ctx, err := newTestGeneratorContext()
	require.NoError(t, err)

	namespaceType := newTestNamespaceType()

	tests := []struct {
		ctx       *generator.Context
		root      *types.Type
		member    types.Member
		hasParent bool
		embedded  *types.Type
		wantErr   error
		want      string
	}{
		{
			// nil pointer t
			ctx:     ctx,
			root:    nil,
			wantErr: errors.New("nil pointer"),
		},
		{
			// test default setter
			ctx:       ctx,
			root:      namespaceType,
			member:    namespaceType.Members[0].Type.Members[objectMetaMember].Type.Members[objectMetaNameMember], // Namespace -> EmbeddedNamespace -> ObjectMeta -> Name (string)
			hasParent: true,
			wantErr:   nil,
			want: `// WithName is an autogenerated setter function.
func (o *Namespace) WithName(in string) *Namespace {
o.Namespace.Name = in
return o
}

`,
		},
		{
			// test default setter, no parent
			ctx:       ctx,
			root:      namespaceType,
			member:    namespaceType.Members[0].Type.Members[objectMetaMember].Type.Members[objectMetaNameMember], // Namespace -> EmbeddedNamespace -> ObjectMeta -> Name (string)
			hasParent: false,
			wantErr:   nil,
			want: `// WithName is an autogenerated setter function.
func (o *Namespace) WithName(in string) *Namespace {
o.Name = in
return o
}

`,
		},
		{
			// test map setter
			ctx:       ctx,
			root:      namespaceType,
			member:    namespaceType.Members[0].Type.Members[objectMetaMember].Type.Members[objectMetaLabelsMember], // Namespace -> EmbeddedNamespace -> ObjectMeta -> Lables (map[string]string)
			hasParent: true,
			wantErr:   nil,
			want: `// WithLabels is an autogenerated setter function.
func (o *Namespace) WithLabels(in map[string]string) *Namespace {
o.Namespace.Labels = mergeMapStringString(o.Namespace.Labels, in)
return o
}

`,
		},
		{
			// test slice embedded type
			ctx:       ctx,
			root:      namespaceType,
			member:    namespaceType.Members[0].Type.Members[objectMetaMember].Type.Members[objectMetaDeploymentsMember], // Namespace -> EmbeddedNamespace -> ObjectMeta -> Deployments ([]Deployment)
			hasParent: true,
			embedded:  namespaceType.Members[0].Type.Members[objectMetaMember].Type.Members[objectMetaDeploymentsMember].Type.Elem.Members[0].Type, // Embedded k8s.io/api/apps/v1.Deployment
			wantErr:   nil,
			want: `// AppendDeployments is an autogenerated setter function.
func (o *Namespace) AppendDeployments(in ...Deployment) *Namespace {
for _, elem := range in {
o.Namespace.Deployments = append(o.Namespace.Deployments, elem.Deployment)
}
return o
}

`,
		},
		{
			// test slice builtin type
			ctx:       ctx,
			root:      namespaceType,
			member:    namespaceType.Members[0].Type.Members[objectMetaMember].Type.Members[objectMetaStringsMember], // Namespace -> EmbeddedNamespace -> ObjectMeta -> Strings ([]string)
			hasParent: true,
			wantErr:   nil,
			want: `// AppendStrings is an autogenerated setter function.
func (o *Namespace) AppendStrings(in ...string) *Namespace {
o.Namespace.Strings = append(o.Namespace.Strings, in...)
return o
}

`,
		},
		{
			// test struct embedded type
			ctx:       ctx,
			root:      namespaceType,
			member:    namespaceType.Members[0].Type.Members[objectMetaMember].Type.Members[objectMetaDeploymentMember], // Namespace -> EmbeddedNamespace -> ObjectMeta -> Deployment (Deployment struct)
			hasParent: true,
			embedded:  namespaceType.Members[0].Type.Members[objectMetaMember].Type.Members[objectMetaDeploymentMember].Type.Members[0].Type, // Embedded k8s.io/api/apps/v1.Deployment
			wantErr:   nil,
			want: `// WithDeployment is an autogenerated setter function.
func (o *Namespace) WithDeployment(in Deployment) *Namespace {
o.Namespace.Deployment = in.Deployment
return o
}

`,
		},
		{
			// test struct type, no embedded
			ctx:       ctx,
			root:      namespaceType,
			member:    namespaceType.Members[0].Type.Members[objectMetaMember].Type.Members[objectMetaSimpleStructMember], // Namespace -> EmbeddedNamespace -> ObjectMeta -> SimpleStruct (struct)
			hasParent: true,
			wantErr:   nil,
			want: `// WithSimpleStruct is an autogenerated setter function.
func (o *Namespace) WithSimpleStruct(in SimpleStruct) *Namespace {
o.Namespace.SimpleStruct = in
return o
}

`,
		},
		{
			// test pointer to builtin
			ctx:       ctx,
			root:      namespaceType,
			member:    namespaceType.Members[0].Type.Members[objectMetaMember].Type.Members[objectMetaIntPtrMember], // Namespace -> EmbeddedNamespace -> ObjectMeta -> IntPtr (*int)
			hasParent: true,
			wantErr:   nil,
			want: `// WithIntPtr is an autogenerated setter function.
func (o *Namespace) WithIntPtr(in int) *Namespace {
o.Namespace.IntPtr = &in
return o
}

`,
		},
	}

	for _, test := range tests {
		var b bytes.Buffer
		err := GenerateSetter(&b, test.ctx, test.root, test.member, test.hasParent, test.embedded)
		if test.wantErr != nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.want, b.String())
		}
	}
}

type Setter struct {
	Template string
	Args     generator.Args
}

func getBuiltInMembersForType(in *types.Type) []types.Member {
	members := []types.Member{}
	for _, m := range in.Members {
		if m.Type.Kind == types.Builtin {
			members = append(members, m)
		}
	}
	return members
}

func setterPrimitives(in *types.Type, members []types.Member) []*Setter {
	s := []*Setter{}
	for _, m := range members {
		s = append(s, setterPrimitiveForMember(in, m))
	}
	return s
}

func setterPrimitiveForMember(parent *types.Type, m types.Member) *Setter {
	s := &Setter{Args: generator.Args{
		"type":       parent,
		"funcName":   funcName(m),
		"memberName": m.Name,
		"memberType": m.Type,
	}}
	var sb strings.Builder
	sb.WriteString("// $.funcName$ is an autogenerated setter function.\n")
	sb.WriteString("func (o *$.type|raw$) $.funcName$(in $.memberType|raw$) *$.type|raw$ {\n")
	sb.WriteString("o.$.memberName$ = in\n")
	sb.WriteString("return o\n")
	sb.WriteString("}\n\n")
	return s
}

func TestSampleType_SetterPrimitive(t *testing.T) {
	ctx, err := newTestGeneratorContext()
	assert.NoError(t, err)

	sample := newSampleBuilderType(t, "Meta")
	var b bytes.Buffer
	//err = GenerateSetter(&b, ctx, sample, sample.Members[0], false, nil)

	s := setterPrimitives(sample)

	err = writeSnippetWithArgs(&b, ctx, s.Template, s.Args)
	assert.NoError(t, err)
	assert.Contains(t, b.String(), "WithName(in string)")
}
