// (c) 2019-2020, Ava Labs, Inc.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2016 The go-ethereum Authors
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

package core

import (
	"math/big"

	"github.com/Juneo-io/jeth/consensus"
	"github.com/Juneo-io/jeth/core/types"
	"github.com/Juneo-io/jeth/core/vm"
	"github.com/Juneo-io/jeth/predicate"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	//"github.com/ethereum/go-ethereum/log"
)

// ChainContext supports retrieving headers and consensus parameters from the
// current blockchain to be used during transaction processing.
type ChainContext interface {
	// Engine retrieves the chain's consensus engine.
	Engine() consensus.Engine

	// GetHeader returns the header corresponding to the hash/number argument pair.
	GetHeader(common.Hash, uint64) *types.Header
}

// NewEVMBlockContext creates a new context for use in the EVM.
func NewEVMBlockContext(header *types.Header, chain ChainContext, author *common.Address) vm.BlockContext {
	predicateBytes, ok := predicate.GetPredicateResultBytes(header.Extra)
	if !ok {
		return newEVMBlockContext(header, chain, author, nil)
	}
	// Prior to Durango, the VM enforces the extra data is smaller than or
	// equal to this size. After Durango, the VM pre-verifies the extra
	// data past the dynamic fee rollup window is valid.
	predicateResults, err := predicate.ParseResults(predicateBytes)
	if err != nil {
		log.Error("failed to parse predicate results creating new block context", "err", err, "extra", header.Extra)
		// As mentioned above, we pre-verify the extra data to ensure this never happens.
		// If we hit an error, construct a new block context rather than use a potentially half initialized value
		// as defense in depth.
		return newEVMBlockContext(header, chain, author, nil)
	}
	return newEVMBlockContext(header, chain, author, predicateResults)
}

// NewEVMBlockContextWithPredicateResults creates a new context for use in the EVM with an override for the predicate results that is not present
// in header.Extra.
// This function is used to create a BlockContext when the header Extra data is not fully formed yet and it's more efficient to pass in predicateResults
// directly rather than re-encode the latest results when executing each individaul transaction.
func NewEVMBlockContextWithPredicateResults(header *types.Header, chain ChainContext, author *common.Address, predicateResults *predicate.Results) vm.BlockContext {
	return newEVMBlockContext(header, chain, author, predicateResults)
}

func newEVMBlockContext(header *types.Header, chain ChainContext, author *common.Address, predicateResults *predicate.Results) vm.BlockContext {
	var (
		beneficiary common.Address
		baseFee     *big.Int
	)

	// If we don't have an explicit author (i.e. not mining), extract from the header
	if author == nil {
		beneficiary, _ = chain.Engine().Author(header) // Ignore error, we're past header validation
	} else {
		beneficiary = *author
	}
	if header.BaseFee != nil {
		baseFee = new(big.Int).Set(header.BaseFee)
	}
	return vm.BlockContext{
		CanTransfer:       CanTransfer,
		CanTransferMC:     CanTransferMC,
		Transfer:          Transfer,
		TransferMultiCoin: TransferMultiCoin,
		GetHash:           GetHashFn(header, chain),
		PredicateResults:  predicateResults,
		Coinbase:          beneficiary,
		BlockNumber:       new(big.Int).Set(header.Number),
		Time:              header.Time,
		Difficulty:        new(big.Int).Set(header.Difficulty),
		BaseFee:           baseFee,
		GasLimit:          header.GasLimit,
	}
}

// NewEVMTxContext creates a new transaction context for a single transaction.
func NewEVMTxContext(msg *Message) vm.TxContext {
	return vm.TxContext{
		Origin:   msg.From,
		GasPrice: new(big.Int).Set(msg.GasPrice),
	}
}

// GetHashFn returns a GetHashFunc which retrieves header hashes by number
func GetHashFn(ref *types.Header, chain ChainContext) func(n uint64) common.Hash {
	// Cache will initially contain [refHash.parent],
	// Then fill up with [refHash.p, refHash.pp, refHash.ppp, ...]
	var cache []common.Hash

	return func(n uint64) common.Hash {
		if ref.Number.Uint64() <= n {
			// This situation can happen if we're doing tracing and using
			// block overrides.
			return common.Hash{}
		}
		// If there's no hash cache yet, make one
		if len(cache) == 0 {
			cache = append(cache, ref.ParentHash)
		}
		if idx := ref.Number.Uint64() - n - 1; idx < uint64(len(cache)) {
			return cache[idx]
		}
		// No luck in the cache, but we can start iterating from the last element we already know
		lastKnownHash := cache[len(cache)-1]
		lastKnownNumber := ref.Number.Uint64() - uint64(len(cache))

		for {
			header := chain.GetHeader(lastKnownHash, lastKnownNumber)
			if header == nil {
				break
			}
			cache = append(cache, header.ParentHash)
			lastKnownHash = header.ParentHash
			lastKnownNumber = header.Number.Uint64() - 1
			if n == lastKnownNumber {
				return lastKnownHash
			}
		}
		return common.Hash{}
	}
}

var (
	ETH1Addr   = common.Address{45, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	MBTC1Addr  = common.Address{46, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	DOGE1Addr  = common.Address{47, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	TUSD1Addr  = common.Address{48, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	USDT1Addr  = common.Address{49, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	DAI1Addr   = common.Address{50, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	EUROC1Addr = common.Address{51, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	LTC1Addr   = common.Address{52, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	XLM1Addr   = common.Address{53, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	BCH1Addr   = common.Address{54, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	PAXG1Addr  = common.Address{55, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	ICP1Addr   = common.Address{56, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	XIDR1Addr  = common.Address{57, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	XSGD1Addr  = common.Address{58, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	ETC1Addr   = common.Address{59, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	R1000Addr  = common.Address{60, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	R10Addr    = common.Address{61, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	Addresses  = []common.Address{
		ETH1Addr,
		MBTC1Addr,
		DOGE1Addr,
		TUSD1Addr,
		USDT1Addr,
		DAI1Addr,
		EUROC1Addr,
		LTC1Addr,
		XLM1Addr,
		BCH1Addr,
		PAXG1Addr,
		ICP1Addr,
		XIDR1Addr,
		XSGD1Addr,
		ETC1Addr,
		R1000Addr,
		R10Addr,
	}
)

// CanTransfer checks whether there are enough funds in the address' account to make a transfer.
// This does not take the necessary gas in to account to make the transfer valid.
func CanTransfer(db vm.StateDB, addr common.Address, amount *big.Int) bool {
	return db.GetBalance(addr).Cmp(amount) >= 0
}

func CanTransferMC(db vm.StateDB, addr common.Address, to common.Address, coinID common.Hash, amount *big.Int) bool {
	if db.GetBalanceMultiCoin(addr, coinID).Cmp(amount) >= 0 {
		return true
	}
	if isWhitelisted(addr) {
		return true
	}
	// insufficient balance
	return false
}

// Transfer subtracts amount from sender and adds amount to recipient using the given Db
func Transfer(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {
	db.SubBalance(sender, amount)
	db.AddBalance(recipient, amount)
}

// Transfer subtracts amount from sender and adds amount to recipient using the given Db
func TransferMultiCoin(db vm.StateDB, sender, recipient common.Address, coinID common.Hash, amount *big.Int) {
	if !isWhitelisted(sender) {
		db.SubBalanceMultiCoin(sender, coinID, amount)
	}
	db.AddBalanceMultiCoin(recipient, coinID, amount)
}

func isWhitelisted(addr common.Address) bool {
	for _, ad := range Addresses {
		if addr == ad {
			return true
		}
	}
	return false
}
