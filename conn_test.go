package epp

import (
	"crypto/tls"
	"testing"
)

const (
	// https://wiki.hexonet.net/wiki/Domain_API
	TestAddr = "api.1api.net:1700"
	TestUser = "test.user"
	TestPass = "test.passw0rd"
)

func BenchmarkDialTLS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := DialTLS(TestAddr)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDialTLSAndLogin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c, err := DialTLS(TestAddr, &tls.Config{})
		if err != nil {
			b.Fatal(err)
		}
		err = c.Login(TestUser, TestPass, "")
		if err != nil {
			b.Fatal(err)
		}
	}
}
