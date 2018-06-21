package eth

import (
	"time"

	gorp "gopkg.in/gorp.v1"
)

const RPC_TIMEOUT = time.Duration(5 * time.Second)

type EthManager struct {
	ethJsonRPC     string
	pollIntervalMS time.Duration
	dbMap          *gorp.DbMap
}

func NewManager(ethJsonRPC string, pollInterval int, dbMap *gorp.DbMap) *EthManager {
	return &EthManager{
		ethJsonRPC:     ethJsonRPC,
		pollIntervalMS: time.Duration(pollInterval * 1000),
		dbMap:          dbMap,
	}
}
