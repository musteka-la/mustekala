package devp2p

import (
	gethlog "github.com/ethereum/go-ethereum/log"
)

// we implement here the interfaces of Logger and Handler
// from github.com/ethereum/go-ethereum
// then we give it to the p2p server as a parameter,
// giving us the ability to log this library.

// this bit does a static check of the interface implementation.
// very useful to tell you at once if your impl is working or not.
var _ gethlog.Logger = (*p2pLibLogger)(nil)

// p2pLibLogger complies with go-ethereum/log interface
type p2pLibLogger struct {
	mgr *Manager
}

// createCustomP2PLibLogger is the function that the manager uses.
// we avoided New, as there is an interface one called that way.
func createCustomP2PLibLogger(m *Manager) *p2pLibLogger {
	return &p2pLibLogger{
		mgr: m,
	}
}

// New complies with go-ethereum/log interface
func (l *p2pLibLogger) New(ctx ...interface{}) gethlog.Logger {
	return &p2pLibLogger{
		mgr: l.mgr,
	}
}

// GetHandler complies with go-ethereum/log interface
func (l *p2pLibLogger) GetHandler() gethlog.Handler {
	return &p2pLibHandler{}
}

// SetHandler complies with go-ethereum/log interface
func (l *p2pLibLogger) SetHandler(h gethlog.Handler) {}

// Trace complies with go-ethereum/log interface and will send the received input to our catchall function
func (l *p2pLibLogger) Trace(msg string, ctx ...interface{}) {
	l.mgr.p2pLibLoggerCatchAll("trace", msg, ctx...)
}

// Debug complies with go-ethereum/log interface and will send the received input to our catchall function
func (l *p2pLibLogger) Debug(msg string, ctx ...interface{}) {
	l.mgr.p2pLibLoggerCatchAll("debug", msg, ctx...)
}

// Info complies with go-ethereum/log interface and will send the received input to our catchall function
func (l *p2pLibLogger) Info(msg string, ctx ...interface{}) {
	l.mgr.p2pLibLoggerCatchAll("info", msg, ctx...)
}

// Warn complies with go-ethereum/log interface and will send the received input to our catchall function
func (l *p2pLibLogger) Warn(msg string, ctx ...interface{}) {
	l.mgr.p2pLibLoggerCatchAll("warn", msg, ctx...)
}

// Error complies with go-ethereum/log interface and will send the received input to our catchall function
func (l *p2pLibLogger) Error(msg string, ctx ...interface{}) {
	l.mgr.p2pLibLoggerCatchAll("error", msg, ctx...)
}

// Crit complies with go-ethereum/log interface and will send the received input to our catchall function
func (l *p2pLibLogger) Crit(msg string, ctx ...interface{}) {
	l.mgr.p2pLibLoggerCatchAll("crit", msg, ctx...)
}

// p2pLibHandler complies with go-ethereum/log interface
type p2pLibHandler struct{}

// Logw complies ith go-ethereum/log interface
func (h *p2pLibHandler) Log(r *gethlog.Record) error {
	return nil
}

///////////////////////////////
// Catch all for this custom logger
///////////////////////////////
// we take it easy with a confy single catch all function with some switches
// and grab what we need.
func (m *Manager) p2pLibLoggerCatchAll(lvl, msg string, ctx ...interface{}) {
	// TODO
	// we may want to use a command flag option to output these logs or not to our lib logger
	log.Debugf("p2p Lib Logger: LEVEL: %v MSG: %v CTX: %v", lvl, msg, ctx)

	switch {
	case lvl == "trace":
		switch {
		case msg == "New dial task":

		}
	}
}
