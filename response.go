package epp

type response_ struct {
	result   Result
	greeting Greeting
	dcr      DomainCheckResponse
}

var scanResponse = NewScanner()
