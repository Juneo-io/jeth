// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"math/big"
	"sync"
	"time"

	"github.com/ava-labs/coreth/params"
	"github.com/ava-labs/coreth/utils"
)

type gasPriceUpdater struct {
	setter       gasPriceSetter
	chainConfig  *params.ChainConfig
	shutdownChan <-chan struct{}

	wg *sync.WaitGroup
}

type gasPriceSetter interface {
	SetGasPrice(price *big.Int)
	SetMinFee(price *big.Int)
}

// handleGasPriceUpdates creates and runs an instance of
func (vm *VM) handleGasPriceUpdates() {
	gpu := &gasPriceUpdater{
		setter:       vm.txPool,
		chainConfig:  vm.chainConfig,
		shutdownChan: vm.shutdownChan,
		wg:           &vm.shutdownWg,
	}

	gpu.start()
}

// start handles the appropriate gas price and minimum fee updates required by [gpu.chainConfig]
func (gpu *gasPriceUpdater) start() {
	// Sets the initial gas price to the launch minimum gas price
	gpu.setter.SetGasPrice(big.NewInt(params.LaunchMinGasPrice))

	// Updates to the minimum gas price as of ApricotPhase1 if it's already in effect or starts a goroutine to enable it at the correct time
	if disabled := gpu.handleUpdate(gpu.setter.SetGasPrice, gpu.chainConfig.ApricotPhase1BlockTimestamp, big.NewInt(params.ApricotPhase1MinGasPrice)); disabled {
		return
	}
	// Updates to the minimum gas price as of ApricotPhase3 if it's already in effect or starts a goroutine to enable it at the correct time
	if disabled := gpu.handleUpdate(gpu.setter.SetGasPrice, gpu.chainConfig.ApricotPhase3BlockTimestamp, big.NewInt(0)); disabled {
		return
	}
	if disabled := gpu.handleUpdate(gpu.setter.SetMinFee, gpu.chainConfig.ApricotPhase3BlockTimestamp, big.NewInt(params.ApricotPhase3MinBaseFee)); disabled {
		return
	}
	// Updates to the minimum gas price as of ApricotPhase4 if it's already in effect or starts a goroutine to enable it at the correct time
	switch {
	case gpu.chainConfig.ChainID.Cmp(params.JUNEChainID) == 0:
		gpu.handleUpdate(gpu.setter.SetMinFee, gpu.chainConfig.ApricotPhase4BlockTimestamp, big.NewInt(params.JUNEChainMinBaseFee))
	case gpu.chainConfig.ChainID.Cmp(params.ETH1ChainID) == 0:
		gpu.handleUpdate(gpu.setter.SetMinFee, gpu.chainConfig.ApricotPhase4BlockTimestamp, big.NewInt(params.ETH1ChainMinBaseFee))
	case gpu.chainConfig.ChainID.Cmp(params.MBTC1ChainID) == 0:
		gpu.handleUpdate(gpu.setter.SetMinFee, gpu.chainConfig.ApricotPhase4BlockTimestamp, big.NewInt(params.MBTC1ChainMinBaseFee))
	case gpu.chainConfig.ChainID.Cmp(params.DOGE1ChainID) == 0:
		gpu.handleUpdate(gpu.setter.SetMinFee, gpu.chainConfig.ApricotPhase4BlockTimestamp, big.NewInt(params.DOGE1ChainMinBaseFee))
	case gpu.chainConfig.ChainID.Cmp(params.TUSD1ChainID) == 0:
		gpu.handleUpdate(gpu.setter.SetMinFee, gpu.chainConfig.ApricotPhase4BlockTimestamp, big.NewInt(params.TUSD1ChainMinBaseFee))
	case gpu.chainConfig.ChainID.Cmp(params.USDT1ChainID) == 0:
		gpu.handleUpdate(gpu.setter.SetMinFee, gpu.chainConfig.ApricotPhase4BlockTimestamp, big.NewInt(params.USDT1ChainMinBaseFee))
	case gpu.chainConfig.ChainID.Cmp(params.DAI1ChainID) == 0:
		gpu.handleUpdate(gpu.setter.SetMinFee, gpu.chainConfig.ApricotPhase4BlockTimestamp, big.NewInt(params.DAI1ChainMinBaseFee))
	case gpu.chainConfig.ChainID.Cmp(params.EUROC1ChainID) == 0:
		gpu.handleUpdate(gpu.setter.SetMinFee, gpu.chainConfig.ApricotPhase4BlockTimestamp, big.NewInt(params.EUROC1ChainMinBaseFee))
	case gpu.chainConfig.ChainID.Cmp(params.LTC1ChainID) == 0:
		gpu.handleUpdate(gpu.setter.SetMinFee, gpu.chainConfig.ApricotPhase4BlockTimestamp, big.NewInt(params.LTC1ChainMinBaseFee))
	case gpu.chainConfig.ChainID.Cmp(params.XLM1ChainID) == 0:
		gpu.handleUpdate(gpu.setter.SetMinFee, gpu.chainConfig.ApricotPhase4BlockTimestamp, big.NewInt(params.XLM1ChainMinBaseFee))
	case gpu.chainConfig.ChainID.Cmp(params.BCH1ChainID) == 0:
		gpu.handleUpdate(gpu.setter.SetMinFee, gpu.chainConfig.ApricotPhase4BlockTimestamp, big.NewInt(params.BCH1ChainMinBaseFee))
	case gpu.chainConfig.ChainID.Cmp(params.PAXG1ChainID) == 0:
		gpu.handleUpdate(gpu.setter.SetMinFee, gpu.chainConfig.ApricotPhase4BlockTimestamp, big.NewInt(params.PAXG1ChainMinBaseFee))
	case gpu.chainConfig.ChainID.Cmp(params.ICP1ChainID) == 0:
		gpu.handleUpdate(gpu.setter.SetMinFee, gpu.chainConfig.ApricotPhase4BlockTimestamp, big.NewInt(params.ICP1ChainMinBaseFee))
	case gpu.chainConfig.ChainID.Cmp(params.XIDR1ChainID) == 0:
		gpu.handleUpdate(gpu.setter.SetMinFee, gpu.chainConfig.ApricotPhase4BlockTimestamp, big.NewInt(params.XIDR1ChainMinBaseFee))
	case gpu.chainConfig.ChainID.Cmp(params.XSGD1ChainID) == 0:
		gpu.handleUpdate(gpu.setter.SetMinFee, gpu.chainConfig.ApricotPhase4BlockTimestamp, big.NewInt(params.XSGD1ChainMinBaseFee))
	case gpu.chainConfig.ChainID.Cmp(params.ETC1ChainID) == 0:
		gpu.handleUpdate(gpu.setter.SetMinFee, gpu.chainConfig.ApricotPhase4BlockTimestamp, big.NewInt(params.ETC1ChainMinBaseFee))
	case gpu.chainConfig.ChainID.Cmp(params.R1000ChainID) == 0:
		gpu.handleUpdate(gpu.setter.SetMinFee, gpu.chainConfig.ApricotPhase4BlockTimestamp, big.NewInt(params.R1000ChainMinBaseFee))
	case gpu.chainConfig.ChainID.Cmp(params.R10ChainID) == 0:
		gpu.handleUpdate(gpu.setter.SetMinFee, gpu.chainConfig.ApricotPhase4BlockTimestamp, big.NewInt(params.R10ChainMinBaseFee))
	default:
		gpu.handleUpdate(gpu.setter.SetMinFee, gpu.chainConfig.ApricotPhase4BlockTimestamp, big.NewInt(params.ApricotPhase4MinBaseFee))
	}
}

// handleUpdate handles calling update(price) at the appropriate time based on
// the value of [timestamp].
// 1) If [timestamp] is nil, update is never called
// 2) If [timestamp] has already passed, update is called immediately
// 3) [timestamp] is some time in the future, starts a goroutine that will call update(price) at the time
// given by [timestamp].
func (gpu *gasPriceUpdater) handleUpdate(update func(price *big.Int), timestamp *uint64, price *big.Int) bool {
	if timestamp == nil {
		return true
	}

	currentTime := time.Now()
	upgradeTime := utils.Uint64ToTime(timestamp)
	if currentTime.After(upgradeTime) {
		update(price)
	} else {
		gpu.wg.Add(1)
		go gpu.updatePrice(update, time.Until(upgradeTime), price)
	}
	return false
}

// updatePrice calls update(updatedPrice) after waiting for [duration] or shuts down early
// if the [shutdownChan] is closed.
func (gpu *gasPriceUpdater) updatePrice(update func(price *big.Int), duration time.Duration, updatedPrice *big.Int) {
	defer gpu.wg.Done()
	select {
	case <-time.After(duration):
		update(updatedPrice)
	case <-gpu.shutdownChan:
	}
}
