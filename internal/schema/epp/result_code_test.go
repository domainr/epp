package epp_test

import (
	"strings"
	"testing"

	"github.com/domainr/epp/internal/schema/epp"
)

func TestResultCodeMessage(t *testing.T) {
	for c := epp.ResultCodeMin; c <= epp.ResultCodeMax; c++ {
		got := c.Message()
		want := epp.Message{Lang: "en", Value: c.String()}
		if got != want {
			t.Errorf("epp.ResultCode(%04d).Message() = %v, want %v", c, got, want)
		}
	}
}

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

func TestResultCodeError(t *testing.T) {
	for c := epp.ResultCodeMin; c <= epp.ResultCodeMax; c++ {
		gotErr := c.Error() != ""
		wantErr := c.IsError()
		if gotErr != wantErr {
			var want string
			if wantErr {
				want = c.String()
			}
			t.Errorf("epp.ResultCode(%04d).Error() = %q, want %q", c, c.Error(), want)
		}
	}
}

func TestResultCodeString(t *testing.T) {
	var known int
	for c := epp.ResultCodeMin; c <= epp.ResultCodeMax; c++ {
		if !strings.HasPrefix(c.String(), "Status code ") {
			known++
		}
	}
	if known != epp.KnownResultCodes {
		t.Errorf("ResultCode values with known string values: %d, want %d", known, epp.KnownResultCodes)
	}
}
