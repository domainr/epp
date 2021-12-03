# EPP for Go

[![build status](https://img.shields.io/github/workflow/status/domainr/epp/Go.svg)](https://github.com/domainr/epp/actions)
[![pkg.go.dev](https://img.shields.io/badge/docs-pkg.go.dev-blue.svg)](https://pkg.go.dev/github.com/domainr/epp)

EPP ([Extensible Provisioning Protocol](https://tools.ietf.org/html/rfc5730)) client for
[Go](https://golang.org/). Extracted from and in production use at [Domainr](https://domainr.com/).

**Note:** This library is currently under development. Its API is subject to breaking changes at any time.

## Installation

`go get github.com/domainr/epp`

## Usage

```go
tconn, err := tls.Dial("tcp", "epp.example.com:700", nil)
if err != nil {
	return err
}

conn, err := epp.NewConn(tconn)
if err != nil {
	return err
}

err = conn.Login(user, password, "")
if err != nil {
	return err
}

dcr, err := conn.CheckDomain("google.com")
if err != nil {
	return err
}
for _, r := range dcr.Checks {
	// ...
}
```

## Todo

- [X] Tests
- [ ] Commands other than `Check`

## Author

© 2021 nb.io LLC
