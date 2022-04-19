package builder

import (
	"k8s.io/gengo/types"

	"k8s.io/gengo/generator"
)

// TODO: this will be refactored out in another PR
type DefaultBuilderFactory struct {
	generator.DefaultGen
}

func (d *DefaultBuilderFactory) NewBuilder(outputFileBaseName string, pkg *types.Package, tagValue string) generator.Generator {
	return d
}
