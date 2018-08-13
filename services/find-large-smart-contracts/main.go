package main

import (
	"bytes"
	"flag"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"

	"github.com/metamask/mustekala/services/find-large-smart-contracts/lib"
)

var NUMBER_OF_NODES = 10000

func main() {
	var (
		blockNumber uint64
		dbFilePath  string
	)

	// command line options
	flag.Uint64Var(&blockNumber, "block-number", 6035070, "Canonical number of the block state to import")
	flag.StringVar(&dbFilePath, "geth-db-filepath", "", "Path to the Go-Ethereum Database")
	flag.Parse()

	// welcome message
	log.Printf("Find Large Smart Contracts Utility")

	gethDB := lib.NewGethDB(dbFilePath)

	// we need to get the state root from the database
	// will query for its canonical block number
	blockHash := gethDB.GetCanonicalHash(blockNumber)
	headerRLP := gethDB.GetHeaderRLP(blockHash, blockNumber)
	header := new(types.Header)
	if err := rlp.Decode(bytes.NewReader(headerRLP), header); err != nil {
		panic(err)
	}

	log.Printf(
		"For the given canonical block number: %d, found this root: %x",
		blockNumber, header.Root[:6])

	// init the trie
	trieDB := trie.NewDatabase(gethDB.DB())
	tr, err := trie.NewSecure(header.Root, trieDB, 120)
	if err != nil {
		panic(err)
	}

	list := getSliceList()
	for _, path := range list {
		log.Printf("peeking on slice %s", path)

		sliceHead := sliceHeadToKeyBytes(path)
		it := trie.NewSliceIterator(tr, sliceHead)
		addresses := it.GetAddressSmartContractOverSize(NUMBER_OF_NODES)

		for _, add := range addresses {
			log.Printf("\tfound %x", add)
		}
	}

	log.Println("ok. done")
}

func sliceHeadToKeyBytes(input string) []byte {
	if input == "" {
		return nil
	}

	// first we convert each character to its hex counterpart
	output := make([]byte, 0)
	var b byte
	for _, c := range input {
		switch {
		case '0' <= c && c <= '9':
			b = byte(c - '0')
		case 'a' <= c && c <= 'f':
			b = byte(c - 'a' + 10)
		default:
			return nil
		}

		output = append(output, b)
	}

	return output
}

func getSliceList() []string {
	return []string{
		"9f13",
		"1a4f",
		"008f",
		"71fe",
		"0a2f",
		"e4a8",
		"5ee8",
		"8b0d",
		"74f8",
		"bf5f",
		"db9a",
		"72e2",
		"56f3",
		"8e96",
		"1019",
		"66f2",
		"315e",
		"deb4",
		"7c20",
		"987e",
		"0981",
		"35b1",
		"3fee",
		"9409",
		"3750",
		"c216",
		"6a60",
		"a21d",
		"1340",
		"2982",
		"2f3f",
		"b4ca",
		"a148",
		"0042",
		"d424",
		"1ca5",
		"47c2",
		"d8b9",
		"b25c",
		"5a58",
		"8195",
		"93f3",
		"df25",
		"881d",
		"8222",
		"d118",
		"1147",
		"6261",
		"6491",
		"a58e",
		"b6ca",
		"3ec6",
		"d836",
		"0c92",
		"4f93",
		"ad0b",
		"d448",
		"0469",
		"55d0",
		"b879",
		"5683",
		"ec0d",
		"28b6",
		"be6b",
		"4de7",
		"a7dd",
		"4ffa",
		"77c4",
		"b755",
		"2b2c",
		"59c6",
		"63b2",
		"ae4b",
		"cff6",
		"7295",
		"bd00",
		"0e9a",
		"135c",
		"17ac",
		"2cc0",
		"3955",
		"474a",
		"525e",
		"6362",
		"b53e",
		"de47",
		"f184",
		"0977",
		"1504",
		"1ad8",
		"6b26",
		"8b67",
		"8c1f",
		"90e8",
		"9e11",
		"b165",
		"c662",
		"fb05",
		"029b",
		"1023",
		"10ab",
		"2d14",
		"31a8",
		"8943",
		"ec51",
		"1ce9",
		"467a",
		"7a8e",
		"8cc6",
		"91f7",
		"ae08",
		"b6d6",
		"c1cd",
		"dde8",
		"f69a",
		"0b46",
		"0dc6",
		"0e18",
		"28b5",
		"34a4",
		"6dd8",
		"78ed",
		"c0fe",
		"c2d4",
		"d443",
		"ffeb",
		"2485",
		"2fdd",
		"5a67",
		"5d17",
		"69c9",
		"9332",
		"99f0",
		"adce",
		"c53b",
		"dbed",
		"017e",
		"1221",
		"1c54",
		"20d4",
		"314b",
		"3313",
		"35c8",
		"35f4",
		"5796",
		"6bfb",
		"76bd",
		"8e6d",
		"b7e0",
		"d043",
		"d8e8",
		"f94c",
		"1b6a",
		"1e3e",
		"23f9",
		"2b1b",
		"49e1",
		"4b31",
		"6d33",
		"6f0f",
		"6f17",
		"7c87",
		"8408",
		"9324",
		"9b86",
		"ae66",
		"b21f",
		"c720",
		"eb12",
		"0aca",
		"11ac",
		"18b5",
		"1926",
		"19b9",
		"2a61",
		"301d",
		"3c82",
		"4bb4",
		"4dda",
		"500b",
		"71ec",
		"781a",
		"7bc0",
		"8495",
		"9b5a",
		"ad0d",
		"ad70",
		"af7e",
		"b19a",
		"bb91",
		"bf7c",
		"c334",
		"dd33",
		"de19",
		"df83",
		"e81e",
		"f11a",
		"f62e",
		"0f96",
		"1a88",
		"2702",
		"3099",
		"3799",
		"50f0",
		"6861",
		"84d3",
		"8b0b",
		"cb36",
		"cd6c",
		"085f",
		"16e0",
		"217a",
		"43dd",
		"5ebc",
		"86c8",
		"9678",
		"ad29",
		"b779",
		"c386",
		"c50d",
		"d1bc",
		"ea6e",
		"edd5",
		"6eb2",
		"acf1",
		"d528",
		"e181",
		"8a4a",
		"dcb9",
		"ed5e",
		"ed4b",
		"5fb7",
		"600c",
		"8296",
		"a7b2",
		"bc35",
	}
}
