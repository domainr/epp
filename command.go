package epp

import (
	"strconv"
	"sync/atomic"
	"time"
)

var txnID = uint64(time.Now().Unix())

func newTxnID() string {
	return strconv.FormatUint(atomic.AddUint64(&txnID, 1), 16)
}
