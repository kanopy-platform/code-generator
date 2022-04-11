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
		Kind:                      "Struct",
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

func newTestEmbeddedNamespaceType() *types.Type {
	return &types.Type{
		Name: types.Name{
			Package: "k8s.io/api/core/v1",
			Name:    "Namespace",
		},
		Kind:                      "Struct",
		CommentLines:              []string{},
		SecondClosestCommentLines: []string{},
		Members:                   []types.Member{},
		Elem:                      nil,
		Key:                       nil,
		Underlying:                nil,
		Methods:                   map[string]*types.Type{},
		Signature:                 nil,
	}
}
