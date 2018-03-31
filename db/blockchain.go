package db

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/garyburd/redigo/redis"
)

// Blockchain is an interface to be used by the devp2p downloader
// to abstract the Redis DB
type BlockChain struct {
	dbPool *redis.Pool

	insertHeaderChainLock sync.RWMutex
}

func LoadBlockChain(dbPool *redis.Pool) *BlockChain {
	b := &BlockChain{
		dbPool: dbPool,
	}

	b.addByzantiumBlock()
	return b
}

// addByzantiumBlock adds the block header 4,370,000
func (b *BlockChain) addByzantiumBlock() {
	var err error

	conn := b.dbPool.Get()
	defer conn.Close()

	// 0. Add the block header value
	byzantiumBlockHex := "f90207a051bc754831f33817e755039d90af3b20ea1e21905529ddaa03d7ba9f5fc9e6" +
		"6fa01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d4934794f3b9d2c81f2b24b0f" +
		"a0acaaa865b7d9ced5fc2fba0e7a73d3c05829730c750ca483b5a65f8321adb25d8abb9da23a4cbb6473464" +
		"eea0402cdc57e9a4bee851f6a7a568f7090489a186d1ff73c4f060c961b685c4668ba01a5b202e1ab165b5c" +
		"296473c3e644e09984785d9f0af55ec83e52362061258c5b90100080020c30a505011120000104048080c20" +
		"840000000800004080890010008108202000102000200000010000004000a4030c004800001d32022006022" +
		"4280040004250000003810020104c0d3004000850401007002708c01009100800005000220008e601108092" +
		"000850020410000050040000082000401000071044440008086400004a601004a00283810c005702200020a" +
		"2118800180442a0881180e000c2605480008910800228100204540a40040005320000820001488000a281c0" +
		"820111440a100e80e6800c000100840400100140848600000100004801200a0123081800030091401102480" +
		"e00b800100310210000900002080080000088021188870aa357c17a7ead8342ae508366528e8364db378459" +
		"e4420386786978697869a0ea6ff3cb300e92d1aa373dd2e1c5c3031489545b61b358fc478ea2e25ac067cb8" +
		"890cbc4e01e1ffc5a"

	byzantiumBlockBin, err := hex.DecodeString(byzantiumBlockHex)
	if err != nil {
		// You never know
		panic("can't decode hardcoded byzantium block")
	}

	byzantiumHash := "b1fcff633029ee18ab6482b58ff8b6e95dd7c82a954c852157152a7a6d32785e"

	byzantiumNumber := "4370000"

	_, err = conn.Do("SET", byzantiumHash, byzantiumBlockBin)
	if err != nil {
		fmt.Printf("Error setting value in redisDB: %v\n", err)
		os.Exit(1)
	}

	// 0. Add the header to the CHT
	_, err = conn.Do("HSET", "canonical-hash-table", byzantiumNumber, byzantiumHash)
	if err != nil {
		fmt.Errorf("Error setting value in redisDB: %v\n", err)
		os.Exit(1)
	}

	// 1. Set the current header pointer
	_, err = conn.Do(
		"HSET", "current-header-pointer",
		"number", byzantiumNumber,
		"hash", byzantiumHash)
	if err != nil {
		fmt.Errorf("Error setting value in redisDB: %v\n", err)
		os.Exit(1)
	}

	// 2. Add the TD for this byzantium block
	byzantiumTD := "1196768507891266117779"
	err = b.setTD(byzantiumNumber, byzantiumHash, byzantiumTD)
	if err != nil {
		fmt.Errorf("Error setting TD in redisDB: %v\n", err)
		os.Exit(1)
	}

	// You are good to go. Live long and prosper.
}

// setTD modifies the total difficulty for a given block
func (b *BlockChain) setTD(number, hash, td string) error {
	var err error

	conn := b.dbPool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%v:%v", number, hash)
	_, err = conn.Do("HSET", "total-difficulty-table", key, td)
	if err != nil {
		return fmt.Errorf("Error setting value in redisDB: %v", err)
	}

	return nil
}

// HasHeader verifies a header's presence in the local chain (i.e. CHT)
func (b *BlockChain) HasHeader(desiredHash common.Hash, desiredHashNumber uint64) bool {
	var err error

	conn := b.dbPool.Get()
	defer conn.Close()

	// 0. find the hash by its number in the CHT
	number := fmt.Sprintf("%v", desiredHashNumber)
	hashStr, err := redis.String(conn.Do("HGET", "canonical-hash-table", number))
	if err != nil && err != redis.ErrNil {
		fmt.Printf("Error getting value from redisDB: %v\n", err)
		return false
	}
	if hashStr == "" {
		return false
	}

	// 1. check the obtained hash against the desired value
	if hashStr != desiredHash.String()[2:] { // Trim 0x prefix
		return false
	}

	return true
}

// GetHeaderByHash retrieves a header from the local chain.
func (b *BlockChain) GetHeaderByHash(hash common.Hash) *types.Header {
	conn := b.dbPool.Get()
	defer conn.Close()

	headerBin, err := conn.Do("GET", hash.String()[2:])
	if err != nil {
		fmt.Printf("Error getting value from redisDB: %v\n", err)
		return nil
	}

	// Parse your binary
	header := new(types.Header)
	if err := rlp.Decode(bytes.NewReader(headerBin.([]byte)), header); err != nil {
		fmt.Printf("Invalid block header RLP | hash %v err %v\n", hash, err)
		return nil
	}
	return header
}

// CurrentHeader retrieves the head header from the local chain.
func (b *BlockChain) CurrentHeader() *types.Header {
	conn := b.dbPool.Get()
	defer conn.Close()

	hash, err := redis.String(conn.Do("HGET", "current-header-pointer", "hash"))
	if err != nil && err != redis.ErrNil {
		fmt.Printf("Error getting value from redisDB: %v\n", err)
		return nil
	}
	if hash == "" {
		return nil
	}

	return b.GetHeaderByHash(common.HexToHash(hash))
}

// GetTd returns the total difficulty of a local block.
func (b *BlockChain) GetTd(hash common.Hash, number uint64) *big.Int {
	conn := b.dbPool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%v:%v", number, hash.String()[2:])
	tdString, err := redis.String(conn.Do("HGET", "total-difficulty-table", key))
	if err != nil && err != redis.ErrNil {
		fmt.Printf("Error getting value from redisDB: %v\n", err)
		return nil
	}

	if tdString == "" {
		fmt.Printf("td should not be an empty string!\n")
		return nil
	}

	td := new(big.Int)
	td.SetString(tdString, 10)

	return td
}

// InsertHeaderChain inserts a batch of headers into the local chain.
func (b *BlockChain) InsertHeaderChain(headers []*types.Header, checkFreq int) (int, error) {
	var err error

	// here is where we verified the order and continuity of the chain of headers,
	// as well as whether the blocks are validly sealed.
	if err = validateHeaderChain(headers, checkFreq); err != nil {
		return 0, fmt.Errorf("%v", err)
	}

	b.insertHeaderChainLock.Lock()
	defer b.insertHeaderChainLock.Unlock()

	for _, header := range headers {
		if err = b.insertHeader(header); err != nil {
			return 0, fmt.Errorf("failed to insert header %v", err)
		}
	}

	return 0, nil
}

// insertHeader inserts the header into the database and modifies the canonical
// chain if its total difficulty is higher.
func (b *BlockChain) insertHeader(header *types.Header) error {
	conn := b.dbPool.Get()
	defer conn.Close()

	// cache hash and header
	var (
		hash   = header.Hash()
		number = header.Number.Uint64()
	)

	// 0. we calculate the total difficulty we are getting with this new header
	parentLocalTD := b.GetTd(header.ParentHash, number-1)
	if parentLocalTD == nil {
		return fmt.Errorf("can not find the parent header of %x (%v) in the local chain", hash[:8], number)
	}
	externalTD := new(big.Int).Add(header.Difficulty, parentLocalTD)

	// 1. we get the local total difficulty
	currentHeader := b.CurrentHeader()
	localTD := b.GetTd(currentHeader.Hash(), b.CurrentHeader().Number.Uint64())

	// 2. add into the database the header and TD associated to this header
	headerBin, err := rlp.EncodeToBytes(header)
	if err != nil {
		return err
	}
	hashString := fmt.Sprintf("%x", hash)
	_, err = conn.Do("SET", hashString, headerBin)
	if err != nil {
		return fmt.Errorf("Error setting value in redisDB: %v", err)
	}
	err = b.setTD(fmt.Sprintf("%v", number), hashString, externalTD.String())
	if err != nil {
		return fmt.Errorf("Error setting TD in redisDB: %v", err)
	}

	// 3. if the total difficulty is higher than our known, add it to the canonical chain
	if externalTD.Cmp(localTD) > 0 {
		// 3.1. delete any canonical number assignments above the new head
		for i := number + 1; ; i++ {
			h, err := redis.String(conn.Do("HGET", "canonical-hash-table", i))
			if err != nil && err != redis.ErrNil {
				return fmt.Errorf("Error getting value from redisDB: %v", err)
			}
			if h == "" {
				// we stop when we can't find more canonical headers upwards
				break
			}

			// Now we delete the actual value
			_, err = conn.Do("HDEL", "canonical-hash-table", i)
			if err != nil {
				return fmt.Errorf("Error removing value from redisDB: %v", err)
			}
		}

		// 3.2. edit the canonical table
		_, err = conn.Do("HSET", "canonical-hash-table", number, fmt.Sprintf("%x", hash))
		if err != nil {
			return fmt.Errorf("Error setting value in redisDB: %v", err)
		}

		// 3.3. set up the new head
		_, err = conn.Do(
			"HSET", "current-header-pointer",
			"number", number,
			"hash", fmt.Sprintf("%x", hash))
		if err != nil {
			return fmt.Errorf("Error setting value in redisDB: %v", err)
		}
	}

	// you are good to go
	return nil
}
