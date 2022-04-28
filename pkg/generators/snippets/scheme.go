package snippets

import "k8s.io/gengo/generator"

func GenerateAddToScheme(importAliases []string) (string, generator.Args) {
	snippet := `
	func AddToSchemeOrPanic(s *runtime.Scheme){
		$- range .PackageAliases$
			utilruntime.Must($.$.SchemeBuilder.AddToScheme(s))
		$- end$
	}
	`

	args := generator.Args{
		"PackageAliases": importAliases,
	}

	return snippet, args
}
