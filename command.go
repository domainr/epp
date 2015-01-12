package epp

import (
	"strconv"
	"sync/atomic"
	"time"
)

// Command represents an EPP <command> element.
type Command struct {
	// Cmd is any valid EPP command, serializable to XML.
	Cmd interface{}

	// TxnID is a unique transaction ID for this command.
	TxnID string `xml:"clTRID"`
}

// NewCommand returns an initialized command.
func NewCommand(cmd interface{}) *Command {
	return &Command{Cmd: cmd, TxnID: newTxnID()}
}

var txnID = uint64(time.Now().Unix())

func newTxnID() string {
	return strconv.FormatUint(atomic.AddUint64(&txnID, 1), 16)
}
