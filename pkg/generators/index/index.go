package index

import (
	"k8s.io/gengo/types"

	"github.com/kanopy-platform/code-generator/pkg/generators/tags"
	log "github.com/sirupsen/logrus"
)

func BuildPackageIndex(index map[string]*types.Type, pkg *types.Package) map[string]*types.Type {
	for _, t := range pkg.Types {
		if tags.IsTypeEnabled(t) {

			for _, m := range t.Members {
				if m.Embedded {
					if _, ok := index[m.Type.String()]; !ok {
						index[m.Type.String()] = t
						log.Debugf("Indexing %s -> (%s, %s) -- Package -> %s(%s)", m.Type.String(), m.Name, m.Type.Name, pkg.Path, pkg.SourcePath)
					}
				}
			}

			if t.Kind == types.Alias {
				ref := tags.ExtractRef(t)
				log.Debugf("Indexing %s - Kind : %s (%s) - Underling Type: %s, Ref: %s", t.Name, t.Kind, t.Name.String(), t.Underlying.Name.String(), ref)
				index[ref] = t
			}
		}
	}
	return index
}
