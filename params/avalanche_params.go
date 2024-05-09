// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"math/big"

	"github.com/ava-labs/avalanchego/utils/units"
)

// Minimum Gas Price
const (
	// MinGasPrice is the number of nAVAX required per gas unit for a
	// transaction to be valid, measured in wei
	LaunchMinGasPrice        int64 = 470_000_000_000
	ApricotPhase1MinGasPrice int64 = 225_000_000_000

	AvalancheAtomicTxFee = units.MilliAvax

	ApricotPhase1GasLimit uint64 = 8_000_000
	CortinaGasLimit       uint64 = 15_000_000

	// TODO Create FeeConfig to store this data
	// StartFee
	JUNEStartMinBaseFee int64 = 48_000_000_000
	MBTCStartMinBaseFee int64 = 13_000_000_000
	GLDStartMinBaseFee  int64 = 3_000_000_000
	LTCStartMinBaseFee  int64 = 54_000_000_000
	DOGEStartMinBaseFee int64 = 6477_000_000_000
	SGDStartMinBaseFee  int64 = 635_000_000_000
	BCHStartMinBaseFee  int64 = 1_000_000_000
	LINKStartMinBaseFee int64 = 26_000_000_000
	EUROStartMinBaseFee int64 = 433_000_000_000
	USDStartMinBaseFee  int64 = 476_000_000_000
	// CurrentFee
	JUNECurrentMinBaseFee int64 = 144_000_000_000
	MBTCCurrentMinBaseFee int64 = 22_000_000_000
	GLDCurrentMinBaseFee  int64 = 1_000_000_000
	LTCCurrentMinBaseFee  int64 = 17_000_000_000
	DOGECurrentMinBaseFee int64 = 9524_000_000_000
	SGDCurrentMinBaseFee  int64 = 1905_000_000_000
	BCHCurrentMinBaseFee  int64 = 3_000_000_000
	LINKCurrentMinBaseFee int64 = 102_000_000_000
	EUROCurrentMinBaseFee int64 = 1299_000_000_000
	USDCurrentMinBaseFee  int64 = 1429_000_000_000

	ApricotPhase3MinBaseFee               int64  = 75_000_000_000
	ApricotPhase3MaxBaseFee               int64  = 225_000_000_000
	ApricotPhase3InitialBaseFee           int64  = 225_000_000_000
	ApricotPhase3TargetGas                uint64 = 10_000_000
	ApricotPhase4MinBaseFee               int64  = 25_000_000_000
	ApricotPhase4MaxBaseFee               int64  = 1_000_000_000_000
	ApricotPhase4BaseFeeChangeDenominator uint64 = 12
	ApricotPhase5TargetGas                uint64 = 15_000_000
	ApricotPhase5BaseFeeChangeDenominator uint64 = 36

	DynamicFeeExtraDataSize        = 80
	RollupWindow            uint64 = 10

	// The base cost to charge per atomic transaction. Added in Apricot Phase 5.
	AtomicTxBaseCost uint64 = 10_000
)

// The atomic gas limit specifies the maximum amount of gas that can be consumed by the atomic
// transactions included in a block and is enforced as of ApricotPhase5. Prior to ApricotPhase5,
// a block included a single atomic transaction. As of ApricotPhase5, each block can include a set
// of atomic transactions where the cumulative atomic gas consumed is capped by the atomic gas limit,
// similar to the block gas limit.
//
// This value must always remain <= MaxUint64.
var AtomicGasLimit *big.Int = big.NewInt(100_000)
