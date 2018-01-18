package devp2p

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
//
// 	"github.com/ethereum/go-ethereum/p2p"
// 	"github.com/ethereum/go-ethereum/p2p/discover"
// 	logging "github.com/ipfs/go-log"
//
// 	"github.com/hermanjunge/devp2p-concept/bridge"
// )

type peerData struct {
  responsive bool
  networkId int
  GenesisBlock string
  isByzantine bool
}

type myMetrics struct {
  data map[string] peerData
  lock sync.RWMutex
}

// metrics := make(map[string] peerData)

func (this *myMetrics) addMetric(peerId string, data *peerData) {
  this.lock.Lock()
  defer this.lock.Unlock()

  this.data[peerId]
}
