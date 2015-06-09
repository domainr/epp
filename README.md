# EPP for Go

[![build status](https://img.shields.io/circleci/project/domainr/epp/master.svg)](https://circleci.com/gh/domainr/epp)
[![godoc](http://img.shields.io/badge/docs-GoDoc-blue.svg)](https://godoc.org/github.com/domainr/epp)

EPP ([Extensible Provisioning Protocol](https://tools.ietf.org/html/rfc5730)) client for
[Go](https://golang.org/). Extracted from and in production use at [Domainr](https://domainr.com/).

**Note:** This library is currently under development, and its interface is subject to breaking changes at any time.

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

dc, err := conn.CheckDomain("google.com")
if err != nil {
	return err
}
for _, r := range dc.Results {
	// ...
}
```

## Todo

- [X] Tests
- [ ] Commands other than `Check`

## Author

© 2015 nb.io, LLC.
