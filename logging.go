package epp

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

// DebugLogger is an io.Writer. Set to enable logging of EPP message XML.
var DebugLogger io.Writer

func logXML(pfx string, p []byte) {
	if DebugLogger == nil {
		return
	}

	var b bytes.Buffer
	enc := xml.NewEncoder(&b)
	enc.Indent("", "\t")

	dec := xml.NewDecoder(bytes.NewReader(p))
	var t xml.Token
	var err error
	for {
		t, err = dec.RawToken()
		if err == io.EOF {
			err = enc.Flush()
			break
		}
		if err != nil {
			break
		}
		err = enc.EncodeToken(t)
		if err != nil {
			break
		}
	}
	if err != nil {
		fmt.Fprintf(DebugLogger, "Indentation error. Raw XML: %s\n%s\n\n", pfx, string(p))
		return
	}

	fmt.Fprintf(DebugLogger, "%s (pretty-printed)\n", pfx)
	io.Copy(DebugLogger, &b)
	fmt.Fprint(DebugLogger, "\n\n")
}
