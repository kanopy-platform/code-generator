package meta

// mock TypeMeta
type TypeMeta struct {
	Kind string
}

// mock ObjectMeta
type ObjectMeta struct {
	Name       string
	Labels     map[string]string
	Finalizers []string
	MockStruct MockStruct
	IntPtr     *int
}

type MockStruct struct {
}
