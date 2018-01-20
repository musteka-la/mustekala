package devp2p

import (
	"crypto/ecdsa"
	"net"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
)

// newServer prepares the devp2p server using the values set in configuration,
// and passed as parameter.
func (m *Manager) newServer(mgrConfig Config) *p2p.Server {
	dialer := p2p.TCPDialer{&net.Dialer{Timeout: 60 * time.Second}}

	name := getClientName()

	privateKey := getPrivateKey(mgrConfig.PrivateKeyFilePath)

	// the defined custom protocol, contains in its handler all the avenues
	// to let the caller of this package establish
	// a communication with the devp2p network.
	protocols := []p2p.Protocol{
		p2p.Protocol{
			Name:    "eth",
			Version: 63,
			Length:  17,
			Run:     m.protocolHandler,
		},
	}

	serverConfig := p2p.Config{
		BootstrapNodes:  mgrConfig.bootnodes,
		Dialer:          dialer,
		ListenAddr:      ":30303",
		Logger:          m.p2pLibLogger,
		MaxPeers:        1000000,
		MaxPendingPeers: 1000000,
		Name:            name,
		NoDiscovery:     false,
		NodeDatabase:    mgrConfig.NodeDatabasePath,
		PrivateKey:      privateKey,
		Protocols:       protocols,
	}
	server := &p2p.Server{Config: serverConfig}
	log.Debug("new devp2p server configured", "instance:", serverConfig.Name)

	return server
}

// getPrivateKey returns the private key found in the given filePath.
// If no file is found, or the file is empty it will generate a random one.
// Otherwise (and invalid key), it will panic.
func getPrivateKey(filePath string) *ecdsa.PrivateKey {
	// TODO
	// PLACEHOLDER
	privKeyString := "af8bf8bb4c634b8716880aa44e82da72b902144940a56e1fa787505aa513ba46"
	privateKey, _ := crypto.HexToECDSA(privKeyString)
	// PLACEHOLDER

	// TODO
	// if error, zero tolerance (you don't wanna delete that key file)
	// if empty, just randomize one

	return privateKey
}

// getClientName returns the name of the client.
func getClientName() string {
	// TODO
	// PLACEHOLDER
	return "Spanish Flea"
	// PLACEHOLDER
}
