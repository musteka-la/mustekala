package devp2p

import (
	"fmt"
	"strings"

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
	l.mgr.p2pLibLoggerCatcher("trace", msg, ctx...)
}

// Debug complies with go-ethereum/log interface and will send the received input to our catchall function
func (l *p2pLibLogger) Debug(msg string, ctx ...interface{}) {
	l.mgr.p2pLibLoggerCatcher("debug", msg, ctx...)
}

// Info complies with go-ethereum/log interface and will send the received input to our catchall function
func (l *p2pLibLogger) Info(msg string, ctx ...interface{}) {
	l.mgr.p2pLibLoggerCatcher("info", msg, ctx...)
}

// Warn complies with go-ethereum/log interface and will send the received input to our catchall function
func (l *p2pLibLogger) Warn(msg string, ctx ...interface{}) {
	l.mgr.p2pLibLoggerCatcher("warn", msg, ctx...)
}

// Error complies with go-ethereum/log interface and will send the received input to our catchall function
func (l *p2pLibLogger) Error(msg string, ctx ...interface{}) {
	l.mgr.p2pLibLoggerCatcher("error", msg, ctx...)
}

// Crit complies with go-ethereum/log interface and will send the received input to our catchall function
func (l *p2pLibLogger) Crit(msg string, ctx ...interface{}) {
	l.mgr.p2pLibLoggerCatcher("crit", msg, ctx...)
}

// p2pLibHandler complies with go-ethereum/log interface
type p2pLibHandler struct{}

// Log complies ith go-ethereum/log interface
func (h *p2pLibHandler) Log(r *gethlog.Record) error {
	return nil
}

// p2pLibLoggerCatcher here we take it easy with a confy single-catch-all function
// with some switches to be able to grab what we need, and pass it to the peer scrapper,
// which may or not process the received info.
//
// there must be a more elegant way to do this, other than just hacking the logs with an axe.
// for now, this does the job, however.
func (m *Manager) p2pLibLoggerCatcher(lvl, msg string, ctx ...interface{}) {

	// You need to activate the flag `--devp2p-lib-debug` to enjoy these logs.
	if m.config.LibP2PDebug {
		log.Debugf("p2p Lib Logger: LEVEL: %v MSG: %v CTX: %v", lvl, msg, ctx)
	}

	// forget about type casting below
	c := fmt.Sprintf("%v", ctx)
	cs := strings.Split(c, " ")

	// this switch is for when we want to format the input
	// on what's going on in the network status file.
	switch {
	case lvl == "trace":
		switch {
		case msg == "New dial task":
			if c[0:13] == "[task dyndial" {
				m.peerScrapper(c[14:len(c)-1], "00-tcp dialing", "wait")
				return
			}
		case msg == "Dial error":
			if c[0:13] == "[task dyndial" {
				peerid := cs[2] + " " + cs[3]
				details := strings.Join(cs[5:len(cs)], " ")
				details = details[0 : len(details)-1]
				m.peerScrapper(peerid, "19-tcp dial fail", details)
				return
			}
		case msg == "Setting up connection failed":
			if cs[0][1:] != "0000000000000000" {
				peerid := cs[0][1:] + " " + cs[1]
				details := strings.Join(cs[2:len(cs)], " ")
				details = details[0 : len(details)-1]
				m.peerScrapper(peerid, "29-connection setup fail", details)
				return
			}
		}
	}
}
