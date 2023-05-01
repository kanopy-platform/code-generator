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
					if m.Name == "StructField" {
						continue
					}
					if _, ok := index[m.Type.String()]; !ok {
						index[m.Type.String()] = t
						log.Debugf("Indexing %s -> (%s, %s) -- Package -> %s(%s)", m.Type.String(), m.Name, m.Type.Name, pkg.Path, pkg.SourcePath)
					}
				}
			}

			if len(t.Members) == 0 {
				log.Debugf("Type with No members: %s - Kind : %s", t.Name, t.Kind)
				if t.Kind == types.Alias {
					log.Debugf("Indexing %s - Kind : %s (%s)-- Underling Type: %s", t.Name, t.Kind, t.Name.String(), t.Underlying.Name.String())
					log.Debugf("\t %#v", t.Underlying.String())

					index[t.String()] = t
				}
			}
		} else {
			log.Debugf("NOT INDEXING: %s", t.Name)
		}
	}
	return index
}
