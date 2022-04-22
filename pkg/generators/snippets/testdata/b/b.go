package b

type TypeMeta struct {
	Kind       string
	APIVersion string
}

const (
	// Number is the index of the member in the struct
	MemberIndex_ObjectMeta_Name   = 0
	MemberIndex_ObjectMeta_Labels = 1
)

type ObjectMeta struct {
	Name   string
	Labels map[string]string
}
