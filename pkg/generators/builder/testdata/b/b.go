package b

// should generate
type BStruct struct {
	Name string
}

// +kanopy:builder=false
type OptOutStruct struct {
	Name string
}
