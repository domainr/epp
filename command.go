package epp

import (
	"strconv"
	"sync/atomic"
	"time"
)

// Command represents an EPP <command> element.
type Command struct {
	// Individual EPP commands.
	DomainCheck *DomainCheck

	// TxnID is a unique transaction ID for this command.
	TxnID string `xml:"clTRID"`
}

var txnID = uint64(time.Now().Unix())

func newTxnID() string {
	return strconv.FormatUint(atomic.AddUint64(&txnID, 1), 16)
}
