package epp

type response_ struct {
	Result
	Greeting
	DomainCheckResponse
}

var scanResponse = NewScanner()
