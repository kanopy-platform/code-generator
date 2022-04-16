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

func newTestNamespaceType() *types.Type {
	return &types.Type{
		Name: types.Name{
			Package: packageName,
			Name:    "Namespace",
		},
		Kind:                      types.Struct,
		CommentLines:              []string{"+kanopy:builder=true"},
		SecondClosestCommentLines: []string{},
		Members: []types.Member{
			{
				Name:     "Namespace",
				Embedded: true,
				Type:     newTestEmbeddedNamespaceType(),
			},
		},
		Elem:       nil,
		Key:        nil,
		Underlying: nil,
		Methods:    map[string]*types.Type{},
		Signature:  nil,
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
		Kind:                      types.Struct,
		CommentLines:              []string{},
		SecondClosestCommentLines: []string{},
		Members: []types.Member{
			{
				Name:     "ObjectMeta",
				Embedded: true,
				Type:     newTestObjectMetaType(),
			},
		},
		Elem:       nil,
		Key:        nil,
		Underlying: nil,
		Methods:    map[string]*types.Type{},
		Signature:  nil,
	}
}

const (
	// make sure numbers match the index of the Members array below
	objectMetaNameMember   = 0
	objectMetaLabelsMember = 1
	objectMetaIntptrMember = 2
)

func newTestObjectMetaType() *types.Type {
	return &types.Type{
		Name: types.Name{
			Package: "k8s.io/apimachinery/pkg/apis/meta/v1",
			Name:    "ObjectMeta",
		},
		Kind:                      types.Struct,
		CommentLines:              []string{},
		SecondClosestCommentLines: []string{},
		Members: []types.Member{
			{
				Name: "Name",
				Type: newTestStringType(),
			},
			{
				Name: "Labels",
				Type: newMapStringStringType(),
			},
			{
				Name: "Intptr",
				Type: newIntptrType(),
			},
		},
		Elem:       nil,
		Key:        nil,
		Underlying: nil,
		Methods:    map[string]*types.Type{},
		Signature:  nil,
	}
}

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

func newMapStringStringType() *types.Type {
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

func newIntptrType() *types.Type {
	return &types.Type{
		Name: types.Name{
			Package: "",
			Name:    "*int",
		},
		Kind: types.Pointer,
		Elem: newTestIntType(),
	}
}
