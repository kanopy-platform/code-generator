package b

type TypeMeta struct {
	Kind       string
	APIVersion string
}

type AliasOfString string

type ObjectMeta struct {
	Name   string
	Labels map[string]string
}
