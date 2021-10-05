package parse

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	DOLLAR_TOKEN        = "$"
	DOUBLE_DOLLAR_TOKEN = "$$"
	DOLLAR_TOKEN_REPLACEMENT        = "##<&DOLLAR_TOKEN&>##"
	DOUBLE_DOLLAR_TOKEN_REPLACEMENT = "##<&DOUBLE_DOLLAR_TOKEN&>##"
)

// Unexpanded represents and unexpanded but to be expanded key value pair.
type Unexpanded struct {
	ID       string
	RawValue string
	Depends  []string
}

// IsExpanded returns whether the given string, contains any expandable variables
func IsExpanded(val string) bool {
	val = strings.ReplaceAll(val, DOUBLE_DOLLAR_TOKEN, DOUBLE_DOLLAR_TOKEN_REPLACEMENT)
	if strings.Contains(val, DOLLAR_TOKEN) {
		return false
	}
	return true
}

// ParseUnexpanded returns a slice of containing any expandable variables of the given string
func ParseUnexpanded(raw string) (unexpanded []string) {
	unexpanded = []string{}
	val := strings.ReplaceAll(raw, DOUBLE_DOLLAR_TOKEN, DOUBLE_DOLLAR_TOKEN_REPLACEMENT)
	req := make(map[string]struct{})

	if strings.Contains(val, DOLLAR_TOKEN) {
		r := regexp.MustCompile(`\$([\w_]*)`)
		unmatched := r.FindAllString(val, -1)
		for _, u := range unmatched {
			if _, ok := req[u]; ok {
				continue
			}
			unexpanded = append(unexpanded, strings.TrimPrefix(u, "$"))
			req[u] = struct{}{}
		}
	}

	return unexpanded
}

// Expand expands Unexpanded using vars and returns the value
func (u *Unexpanded) Expand(vars map[string]string) string {
	if IsExpanded(u.RawValue) {
		return u.RawValue
	}

	val := strings.ReplaceAll(u.RawValue, DOUBLE_DOLLAR_TOKEN, DOUBLE_DOLLAR_TOKEN_REPLACEMENT)

	for _, dep := range u.Depends {
		if subst, ok := vars[dep]; ok {
			newVal := strings.ReplaceAll(val, fmt.Sprintf("$%s", dep), subst)
			val = newVal
		}
		newVal := strings.ReplaceAll(val, fmt.Sprintf("$%s", dep), "")
		val = newVal
	}

	val = strings.ReplaceAll(val, DOUBLE_DOLLAR_TOKEN_REPLACEMENT, DOUBLE_DOLLAR_TOKEN)
	return val
}
