package epp_test

import (
	"testing"

	"github.com/domainr/epp/internal/schema/epp"
)

func TestResultCodeIsError(t *testing.T) {
	for c := epp.ResultCodeMin; c <= epp.ResultCodeMax; c++ {
		got := c.IsError()
		want := c >= 2000
		if got != want {
			t.Errorf("epp.ResultCode(%04d).IsError() = %t, want %t", c, got, want)
		}
	}
}

func TestResultCodeIsFatal(t *testing.T) {
	for c := epp.ResultCodeMin; c <= epp.ResultCodeMax; c++ {
		got := c.IsFatal()
		want := c >= 2500
		if got != want {
			t.Errorf("epp.ResultCode(%04d).IsFatal() = %t, want %t", c, got, want)
		}
	}
}
