package e

type MockPolicyRule struct {
	Verbs                 []string
	ListOfInts            []int
	AliasType             *AliasToString
	ToggleAliasWithoutRef *AnotherAlias
}

type AliasToString string
type AnotherAlias string
