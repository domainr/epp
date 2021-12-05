package epp

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"sync/atomic"
)

type ID interface {
	ID() string
}

type seqSource struct {
	prefix string
	n      uint64
}

func newSeqSource(prefix string) (*seqSource, error) {
	if prefix == "" {
		var pfx [16]byte
		_, err := rand.Read(pfx[:])
		if err != nil {
			return nil, err
		}
		prefix = hex.EncodeToString(pfx[:])
	}
	return &seqSource{
		prefix: prefix,
	}, nil
}

func (s *seqSource) ID() string {
	return s.prefix + strconv.FormatUint(atomic.AddUint64(&s.n, 1), 10)
}
