package epp

import (
	"context"

	"github.com/domainr/epp/internal/schema/epp"
)

type transaction struct {
	ctx   context.Context
	reply chan reply
}

func newTransaction(ctx context.Context) (transaction, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	return transaction{
		ctx:   ctx,
		reply: make(chan reply),
	}, cancel
}

type reply struct {
	e   *epp.EPP
	err error
}
