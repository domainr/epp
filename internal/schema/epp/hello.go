package epp

type Hello struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 hello,selfclosing"`
}

func (Hello) eppBody() {}
