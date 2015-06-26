package epp

import (
	"fmt"
	"io"
)

// DebugLogger is an io.Writer. Set to enable logging of EPP message XML.
var DebugLogger io.Writer

func logXML(pfx string, xml []byte) {
	if DebugLogger == nil {
		return
	}
	fmt.Printf("%s\n%s\n\n", pfx, string(xml))
}
