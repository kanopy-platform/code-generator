package snippets

import (
	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

func GenerateConstructorForObjectMeta(t *types.Type) (string, generator.Args) {
	args := generator.Args{
		"type": t,
	}

	raw := `// New$.type|raw$ is an autogenerated constructor.
func New$.type|raw$(name string) *$.type|raw$ {
	o := &$.type|raw${}
	o.ObjectMeta.Name = name
	return o
}

`
	return raw, args
}
