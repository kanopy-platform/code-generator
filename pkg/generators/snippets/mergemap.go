package snippets

func GenerateMergeMapStringString() string {
	raw := `// mergeMapStringString creates a new map and loads it from map args
// This function takes at least 2 args. Later map args take precedence.
func mergeMapStringString(m1 map[string]string, mapArgs ...map[string]string) map[string]string {
	outMap := map[string]string{}
	for k, v := range m1 {
		outMap[k] = v
	}

	for _, m := range mapArgs {
		for k, v := range m {
			outMap[k] = v
		}
	}
	return outMap
}

`
	return raw
}
