package expand_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alex-held/expand"
)

func TestNewResolverWithEnviron(t *testing.T) {
	vars := map[string]string{}
	expected := os.ExpandEnv("$HOME")
	r := expand.NewExpanderWithEnvironment(vars)
	resolved, err := r.Expand()
	assert.NoError(t, err)
	assert.Equal(t, expected, resolved["HOME"])
}

func TestResolve(t *testing.T) {
	vars := map[string]string{
		"all":  "$HOME/$$foo/$bar/$baz/SOME_THING",
		"bar":  "bar_value",
		"baz":  "baz_value",
		"HOME": "/user/home", // override $HOME variable for test
	}
	expected := "/user/home/$$foo/bar_value/baz_value/SOME_THING"

	sut := expand.NewExpander(vars)

	resolved, err := sut.Expand()
	assert.NoError(t, err)

	actual := resolved["all"]
	assert.Equal(t, expected, actual)
}
