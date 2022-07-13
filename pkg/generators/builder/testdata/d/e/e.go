package e

type MockPolicyRule struct {
	Verbs      []string
	ListOfInts []int
	AliasType  *AliasToString
}

type AliasToString string
