package snippets

func GenerateVariadicBool() string {
	raw := `// variadicBool selects the first element in the passed in list if non-empty. Otherwise the default return is "true".
func variadicBool(in ...bool) bool {
	if len(in) > 0 {
		return in[0]
	}
	return true
}

`
	return raw
}

func GenerateBoolPointer() string {
	raw := `// boolPointer returns a pointer to a bool.
func boolPointer(in bool) *bool {
	return &in
}

`
	return raw
}
