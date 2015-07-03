package epp

import "encoding/xml"

type response_ struct {
	result  Result
	resData resData
}

func (d *Decoder) decodeResponse(res *response_) error {
	*res = response_{}
	return d.DecodeElements(func(e xml.StartElement) error {
		switch e.Name.Local {
		case "result":
			return d.decodeResult(&res.result)
		case "resData":
			return d.decodeResData(&res.resData)
		}
		return nil
	})
}

type resData struct {
	domainChkData domainChkData
}

func (d *Decoder) decodeResData(rd *resData) error {
	*rd = resData{}
	return d.DecodeWith(func(t xml.Token) error {
		var err error
		switch node := t.(type) {
		case xml.StartElement:
			switch node.Name.Local {
			case "chkData":
				err = d.decodeDomainChkData(&rd.domainChkData)
			}
		}
		return err
	})
}

type domainChkData struct {
	cds []domainCD
}

func (d *Decoder) decodeDomainChkData(v *domainChkData) error {
	*v = domainChkData{}
	return d.DecodeWith(func(t xml.Token) error {
		var err error
		switch node := t.(type) {
		case xml.StartElement:
			switch node.Name.Local {
			case "cd":
				v.cds = append(v.cds, domainCD{})
				err = d.decodeDomainCD(&v.cds[len(v.cds)-1])
			}
		}
		return err
	})
}

type domainCD struct {
	name  string
	avail bool
}

func (d *Decoder) decodeDomainCD(v *domainCD) error {
	*v = domainCD{}
	return d.DecodeWith(func(t xml.Token) error {
		var err error
		// switch node := t.(type) {
		// case xml.StartElement:
		// 	switch node.Name.Local {
		// 	case "cd":
		// 		err = d.decodeDomainCD(&v.cds[len(v.cds)-1])
		// 	}
		// }
		return err
	})
}
