package e

type MockPolicyRule struct {
	Verbs                 []string
	ListOfInts            []int
	AliasType             *AliasToString
	ToggleAliasWithoutRef *AnotherAlias
	privateField          PrivateField
}

type PrivateField struct{}
type AliasToString string
type AnotherAlias string
