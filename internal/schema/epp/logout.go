package epp

type Logout struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 logout,selfclosing"`
}

func (Logout) eppCommand() {}
