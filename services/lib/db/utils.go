// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
package db

// Copyright 2015 The go-ethereum Authors
import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

var (
	bogusEthash      *ethash.Ethash
	bogusChainReader *myChainReader
)

func init() {
	bogusEthash = ethash.New(ethash.Config{})
	bogusChainReader = &myChainReader{}
}

func validateHeaderChain(chain []*types.Header, checkFreq int) error {
	// Do a sanity check that the provided chain is actually ordered and linked
	for i := 1; i < len(chain); i++ {
		if chain[i].Number.Uint64() != chain[i-1].Number.Uint64()+1 || chain[i].ParentHash != chain[i-1].Hash() {
			return fmt.Errorf("non contiguous insert: item %d is #%d [%x…], item %d is #%d [%x…] (parent [%x…])", i-1, chain[i-1].Number,
				chain[i-1].Hash().Bytes()[:4], i, chain[i].Number, chain[i].Hash().Bytes()[:4], chain[i].ParentHash[:4])
		}
	}

	// Iterate over the headers and ensure they all check out
	for _, header := range chain {
		if err := bogusEthash.VerifySeal(bogusChainReader, header); err != nil {
			return fmt.Errorf("Error verifying header %v: %v", header.Number, err)
		}
	}

	return nil
}

///////////
// Bogus Object to comply with the interface
///////////
type myChainReader struct{}

func (b *myChainReader) Config() *params.ChainConfig  { return &params.ChainConfig{} }
func (b *myChainReader) CurrentHeader() *types.Header { return &types.Header{} }
func (b *myChainReader) GetHeader(hash common.Hash, number uint64) *types.Header {
	return &types.Header{}
}
func (b *myChainReader) GetHeaderByNumber(number uint64) *types.Header         { return &types.Header{} }
func (b *myChainReader) GetHeaderByHash(hash common.Hash) *types.Header        { return &types.Header{} }
func (b *myChainReader) GetBlock(hash common.Hash, number uint64) *types.Block { return &types.Block{} }
