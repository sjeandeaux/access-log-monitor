package matcher

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/onsi/gomega/types"
)

//CmpEqual uses cmp.Equal to compare actual with expected.  Equal is strict about
//types when performing comparisons.
//It is an error for both actual and expected to be nil.  Use BeNil() instead.
func CmpEqual(expected interface{}) types.GomegaMatcher {
	return &CmpMatcher{
		Expected: expected,
	}
}

// CmpMatcher compare
type CmpMatcher struct {
	Expected interface{}
}

// Match test if they match
func (matcher *CmpMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil && matcher.Expected == nil {
		return false, fmt.Errorf("refusing to compare <nil> to <nil>")
	}

	return cmp.Equal(actual, matcher.Expected), nil
}

// FailureMessage print the failure
func (matcher *CmpMatcher) FailureMessage(actual interface{}) (message string) {
	return cmp.Diff(actual, matcher.Expected)
}

// NegatedFailureMessage print the failure
func (matcher *CmpMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return cmp.Diff(actual, matcher.Expected)
}
