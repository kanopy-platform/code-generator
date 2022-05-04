package tags

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/gengo/args"
	"k8s.io/gengo/types"
)

func TestExtractCommentTag(t *testing.T) {
	fmtTag := "+%s=%s"
	tests := []struct {
		description string
		tag         string
		comments    []string
		want        string
	}{
		{
			description: "Empty value from no comments",
			tag:         Builder,
			comments:    []string{},
			want:        "",
		},
		{
			description: "Empty value from empty comments",
			tag:         Builder,
			comments:    []string{""},
			want:        "",
		},
		{
			description: "Value from comments",
			tag:         Builder,
			comments:    []string{fmt.Sprintf(fmtTag, Builder, "value")},
			want:        "value",
		},
		{
			description: "Tag with no value",
			tag:         Builder,
			comments:    []string{fmt.Sprintf(fmtTag, Builder, "")},
			want:        "",
		},
		{
			description: "Return first value with multiple values",
			tag:         Builder,
			comments:    []string{fmt.Sprintf(fmtTag, Builder, "value,value2")},
			want:        "value",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.want, Extract(test.comments, test.tag), test.description)
	}
}

func TestTypeEnabled(t *testing.T) {
	assert.True(t, IsTypeEnabled(getTestPackage(t).Types["AType"]))
}

func TestTypeOptOut(t *testing.T) {
	assert.True(t, IsTypeOptedOut(getTestPackage(t).Types["OptType"]))
}

func TestAllPackageTypes(t *testing.T) {
	assert.True(t, IsPackageTagged(getTestPackage(t).Comments))
}

func TestMemberReadyOnly(t *testing.T) {
	testType := getTestPackage(t).Types["MemberComments"]
	assert.False(t, IsMemberReadyOnly(getMemberByName(t, testType, "Name")))
	assert.True(t, IsMemberReadyOnly(getMemberByName(t, testType, "UID")))
}

func getTestPackage(t *testing.T) *types.Package {
	testDir := "./testdata/a"
	d := args.Default()
	d.IncludeTestFiles = true
	d.InputDirs = []string{testDir + ""}
	d.GoHeaderFilePath = filepath.Join(args.DefaultSourceTree())
	b, err := d.NewBuilder()
	assert.NoError(t, err)
	findTypes, err := b.FindTypes()
	assert.NoError(t, err)
	return findTypes[testDir]
}

func getMemberByName(t *testing.T, inputType *types.Type, name string) types.Member {
	for _, m := range inputType.Members {
		if m.Name == name {
			return m
		}
	}

	t.Fatalf("failed to find %q in type %q", name, inputType)
	return types.Member{}
}
