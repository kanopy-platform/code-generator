package a

type NoGeneration struct {
	Name string
}

// +kanopy:builder=true
type AStruct struct {
	Name string
}
