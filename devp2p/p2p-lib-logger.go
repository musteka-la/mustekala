package devp2p

import (
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"
)

// we implement here the interfaces of Logger and Handler
// from github.com/ethereum/go-ethereum
// then we give it to the p2p server as a parameter,
// giving us the ability to log this library.

// this bit does a static check of the interface implementation.
// very useful to tell you at once if your impl is working or not.
var _ gethlog.Logger = (*p2pLibLogger)(nil)

type p2pLibLogger struct{}

func (l *p2pLibLogger) New(ctx ...interface{}) gethlog.Logger {
	return &p2pLibLogger{}
}

func (l *p2pLibLogger) GetHandler() gethlog.Handler {
	return &p2pLibHandler{}
}
func (l *p2pLibLogger) SetHandler(h gethlog.Handler) {}

func (l *p2pLibLogger) Trace(msg string, ctx ...interface{}) { fmt.Printf("Custom Logger: %v\n", msg) }
func (l *p2pLibLogger) Debug(msg string, ctx ...interface{}) { fmt.Printf("Custom Logger: %v\n", msg) }
func (l *p2pLibLogger) Info(msg string, ctx ...interface{})  { fmt.Printf("Custom Logger: %v\n", msg) }
func (l *p2pLibLogger) Warn(msg string, ctx ...interface{})  { fmt.Printf("Custom Logger: %v\n", msg) }
func (l *p2pLibLogger) Error(msg string, ctx ...interface{}) { fmt.Printf("Custom Logger: %v\n", msg) }
func (l *p2pLibLogger) Crit(msg string, ctx ...interface{})  { fmt.Printf("Custom Logger: %v\n", msg) }

type p2pLibHandler struct{}

func (h *p2pLibHandler) Log(r *gethlog.Record) error {
	return nil
}
