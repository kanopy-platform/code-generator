package a

type AStruct struct {
	Meta
	Spec   Spec
	Direct string
}

type Meta struct {
	Name string
}

type Spec struct {
	Complex *AComplex
}

type AComplex struct {
	C1   string
	C2   int
	C3   *int
	Many []string
}
