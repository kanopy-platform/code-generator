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
		comments    []string
		want        string
	}{
		{
			description: "Empty value from no comments",
			comments:    []string{},
			want:        "",
		},
		{
			description: "Empty value from empty comments",
			comments:    []string{""},
			want:        "",
		},
		{
			description: "Value from comments",
			comments:    []string{fmt.Sprintf(fmtTag, Name, "value")},
			want:        "value",
		},
		{
			description: "Tag with no value",
			comments:    []string{fmt.Sprintf(fmtTag, Name, "")},
			want:        "",
		},
		{
			description: "Return first value with multiple values",
			comments:    []string{fmt.Sprintf(fmtTag, Name, "value,value2")},
			want:        "value",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.want, Extract(test.comments), test.description)
	}
}

func TestTypeEnabled(t *testing.T) {
	assert.True(t, IsTypeEnabled(getTestPackage(t).Types["AType"]))
}

func TestTypeOptOut(t *testing.T) {
	assert.True(t, IsTypeOptedOut(getTestPackage(t).Types["OptType"]))
}

func TestAllPackageTypes(t *testing.T) {
	assert.True(t, IsPackageTagged(Extract(getTestPackage(t).Comments)))
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
