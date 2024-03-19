package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncludeMember(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description  string
		dir          string
		typeSelector string
		member       string
		want         bool
	}{
		{
			description:  "should include",
			dir:          "d/e",
			typeSelector: "MockPolicyRule",
			member:       "Verbs",
			want:         true,
		},
		{
			description:  "should not include member labeled Read-only",
			dir:          "c/meta",
			typeSelector: "ObjectMeta",
			member:       "ReadOnlyMember",
			want:         false,
		},
		{
			description:  "should not include member labeled read-only (lowercase)",
			dir:          "c/meta",
			typeSelector: "ObjectMeta",
			member:       "ReadOnlyLowerCase",
			want:         false,
		},
		{
			description:  "should not include private member",
			dir:          "d/e",
			typeSelector: "MockPolicyRule",
			member:       "privateField",
			want:         false,
		},
	}

	for _, test := range tests {
		_, testType := newTestGeneratorType(t, test.dir, test.typeSelector)
		member := getMemberFromType(testType, test.member)
		assert.NotEmpty(t, member, test.description)
		assert.Equal(t, test.want, includeMember(testType, member), test.description)
	}
}

func TestIncludeObjectMetaMember(t *testing.T) {
	t.Parallel()

	_, objectMeta := newTestGeneratorType(t, "c/meta", "ObjectMeta")

	tests := []struct {
		description string
		member      string
		want        bool
	}{
		{
			description: "should include Name",
			member:      "Name",
			want:        true,
		},
		{
			description: "should not include Finalizers",
			member:      "Finalizers",
			want:        false,
		},
	}

	for _, test := range tests {
		member := getMemberFromType(objectMeta, test.member)
		assert.NotEmpty(t, member, test.description)
		assert.Equal(t, test.want, includeObjectMetaMember(member), test.description)
	}
}
