package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKVP_Requires(t *testing.T) {
	tests := []struct {
		name     string
		expected []string
	}{
		{
			name:     "$HOME/$$foo/$bar/$baz/SOME_THING",
			expected: []string{"HOME", "bar", "baz"},
		},
		{
			name:     "$HOME",
			expected: []string{"HOME"},
		},
		{
			name:     "$foo$bar",
			expected: []string{"foo", "bar"},
		},
		{
			name:     "$foo$bar$foo",
			expected: []string{"foo", "bar"},
		},
		{
			name:     "$$$foo",
			expected: []string{"foo"},
		},
		{
			name:     "SOME_THING",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ParseUnexpanded(tt.name)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestExpand_IsExpanded(t *testing.T) {
	type vars struct {
		val string
	}
	tests := []struct {
		name     string
		vars     vars
		expected bool
	}{
		{
			name: "$$ be falsee",
			vars: vars{
				val: "foo/$$bar/baz",
			},
			expected: true,
		},
		{
			name: "$ be false",
			vars: vars{
				val: "foo/$bar/baz",
			},
			expected: false,
		},
		{
			name: "$$ and $ be false",
			vars: vars{
				val: "foo/$bar/baz",
			},
			expected: false,
		},
		{
			name: "no expansions be true",
			vars: vars{
				val: "foo/bar/baz",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := IsExpanded(tt.vars.val)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
