package main

import (
	"errors"
	"strings"
	"testing"
	"testing/quick"
)

func TestFoldErrors(t *testing.T) {
	errPrefix := errors.New("_prefix_")
	errSuffix := errors.New("_suffix_")

	// Note: It is more than twice as slow with this predicate approach... Maybe revert to direct inline calls
	isNotNil := func(err error) bool { return err != nil }
	hasProperPrefix := func(err error) bool { return strings.HasPrefix(err.Error(), errPrefix.Error()) }
	hasProperSuffix := func(err error) bool { return strings.HasSuffix(err.Error(), errSuffix.Error()) }

	errorsGen := func(msg string) bool {
		generatedErr := errors.New(msg)
		errs := []error{errPrefix, generatedErr, errSuffix}

		if foldedErr := foldErrors(errs); isNotNil(foldedErr) && hasProperPrefix(foldedErr) && hasProperSuffix(foldedErr) {
			return true
		}
		return false
	}

	if err := quick.Check(errorsGen, nil); err != nil {
		t.Error(err)
	}
}
