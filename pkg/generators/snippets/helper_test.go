/*
 Helper functions for _test.go files
*/

package snippets

import (
	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
)

const packageName string = "./k8s"

/* Functions for creating the generator.Context with appropriate NameSystems.
 */
func nameSystem() namer.NameSystems {
	return namer.NameSystems{
		"public": namer.NewPublicNamer(1),
		"raw":    namer.NewRawNamer(packageName, nil),
	}
}

func defaultNameSystem() string {
	return "public"
}

func newTestGeneratorContext() (*generator.Context, error) {
	args := args.Default()

	b, err := args.NewBuilder()
	if err != nil {
		return nil, err
	}

	c, err := generator.NewContext(b, nameSystem(), defaultNameSystem())
	if err != nil {
		return nil, err
	}

	return c, nil
}

/* Functions that return mock k8s wrapper types which needs code generation.
 */
func newTestNamespaceType() *types.Type {
	return &types.Type{
		Name: types.Name{
			Package: packageName,
			Name:    "Namespace",
		},
		Kind:         types.Struct,
		CommentLines: []string{"+kanopy:builder=true"},
		Members: []types.Member{
			{
				//Name:     "Namespace",
				Embedded: true,
				Type:     newTestEmbeddedNamespaceType(),
			},
		},
	}
}

const (
	// make sure numbers match the index of the Members array below
	objectMetaMember = 0
)

func newTestEmbeddedNamespaceType() *types.Type {
	return &types.Type{
		Name: types.Name{
			Package: "k8s.io/api/core/v1",
			Name:    "Namespace",
		},
		Kind: types.Struct,
		Members: []types.Member{
			{
				//Name:     "ObjectMeta",
				Embedded: true,
				Type:     newTestObjectMetaType(),
			},
		},
	}
}

func newTestDeploymentType() *types.Type {
	return &types.Type{
		Name: types.Name{
			Package: packageName,
			Name:    "Deployment",
		},
		Kind:                      types.Struct,
		CommentLines:              []string{"+kanopy:builder=true"},
		SecondClosestCommentLines: []string{},
		Members: []types.Member{
			{
				//Name:     "Deployment",
				Embedded: true,
				Type:     newTestEmbeddedDeploymentType(),
			},
		},
		Elem:       nil,
		Key:        nil,
		Underlying: nil,
		Methods:    map[string]*types.Type{},
		Signature:  nil,
	}
}

func newTestEmbeddedDeploymentType() *types.Type {
	return &types.Type{
		Name: types.Name{
			Package: "k8s.io/api/apps/v1",
			Name:    "Deployment",
		},
		Kind:    types.Struct,
		Members: []types.Member{},
	}
}

const (
	// make sure numbers match the index of the Members array below
	objectMetaNameMember         = 0
	objectMetaLabelsMember       = 1
	objectMetaDeploymentsMember  = 2
	objectMetaStringsMember      = 3
	objectMetaDeploymentMember   = 4
	objectMetaSimpleStructMember = 5
	objectMetaIntPtrMember       = 6
)

func newTestObjectMetaType() *types.Type {
	return &types.Type{
		Name: types.Name{
			Package: "k8s.io/apimachinery/pkg/apis/meta/v1",
			Name:    "ObjectMeta",
		},
		Kind: types.Struct,
		Members: []types.Member{
			{
				// Name -> string
				Name: "Name",
				Type: newTestStringType(),
			},
			{
				// Labels -> map[string]string
				Name: "Labels",
				Type: newTestMapStringStringType(),
			},
			{
				// Deployments -> []Deployment
				Name: "Deployments",
				Type: newTestDeploymentSliceType(),
			},
			{
				// Strings -> []string
				Name: "Strings",
				Type: newTestStringSliceType(),
			},
			{
				// Deployment -> Deployment struct
				Name: "Deployment",
				Type: newTestDeploymentType(),
			},
			{
				// SimpleStruct -> struct
				Name: "SimpleStruct",
				Type: newTestSimpleStructType(),
			},
			{
				Name: "IntPtr",
				Type: newTestIntPtrType(),
			},
		},
	}
}

/* Functions that return mock basic and complex types used by the types above
 */
func newTestIntType() *types.Type {
	return &types.Type{
		Name: types.Name{
			Package: "",
			Name:    "int",
		},
		Kind: types.Builtin,
	}
}

func newTestStringType() *types.Type {
	return &types.Type{
		Name: types.Name{
			Package: "",
			Name:    "string",
		},
		Kind: types.Builtin,
	}
}

func newTestStringSliceType() *types.Type {
	return &types.Type{
		Name: types.Name{
			Package: "",
			Name:    "[]string",
		},
		Kind: types.Slice,
		Elem: newTestStringType(),
	}
}

func newTestDeploymentSliceType() *types.Type {
	return &types.Type{
		Name: types.Name{
			Package: "",
			Name:    "[]Deployment",
		},
		Kind: types.Slice,
		Elem: newTestDeploymentType(),
	}
}

func newTestMapStringStringType() *types.Type {
	return &types.Type{
		Name: types.Name{
			Package: "",
			Name:    "map[string]string",
		},
		Kind: types.Map,
		Elem: newTestStringType(),
		Key:  newTestStringType(),
	}
}

func newTestIntPtrType() *types.Type {
	return &types.Type{
		Name: types.Name{
			Package: "",
			Name:    "*int",
		},
		Kind: types.Pointer,
		Elem: newTestIntType(),
	}
}

func newTestSimpleStructType() *types.Type {
	return &types.Type{
		Name: types.Name{
			Package: packageName,
			Name:    "SimpleStruct",
		},
		Kind: types.Struct,
		Members: []types.Member{
			{
				Name: "Description",
				Type: newTestStringType(),
			},
		},
	}
}
