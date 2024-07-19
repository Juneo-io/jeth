# JEth and the Juneo EVM

[Juneo Supernet](https://juneo.com) is a network composed of multiple blockchains.
Each blockchain is an instance of a Virtual Machine (VM), much like an object in an object-oriented language is an instance of a class.
That is, the VM defines the behavior of the blockchain.
JEth (from Juneo Ethereum) is the virtual machine that defines the Juneo EVM.
This chain implements the Ethereum Virtual Machine and supports Solidity smart contracts as well as most other Ethereum client functionality.

## API

The Juneo EVM supports the following API namespaces:

- `eth`
- `personal`
- `txpool`
- `debug`

Only the `eth` namespace is enabled by default.

## Compatibility

The Juneo EVM is compatible with almost all Ethereum tooling, including Metamask, Remix, Truffle and Hardhat.

## Differences Between Juneo JEth and Ethereum

### Atomic Transactions

As a network composed of multiple blockchains, Juneo EVM uses _atomic transactions_ to move assets between chains. JEth modifies the Ethereum block format by adding an _ExtraData_ field, which contains the atomic transactions.

### Block Timing

Blocks are produced asynchronously in Snowman Consensus, so the timing assumptions that apply to Ethereum do not apply to JEth. To support block production in an async environment, a block is permitted to have the same timestamp as its parent. Since there is no general assumption that a block will be produced every 10 seconds, smart contracts built on Juneo EVM should use the block timestamp instead of the block number for their timing assumptions.

A block with a timestamp more than 10 seconds in the future will not be considered valid. However, a block with a timestamp more than 10 seconds in the past will still be considered valid as long as its timestamp is greater than or equal to the timestamp of its parent block.

## Difficulty and Random OpCode

Snowman consensus does not use difficulty in any way, so the difficulty of every block is required to be set to 1. This means that the DIFFICULTY opcode should not be used as a source of randomness.

Additionally, with the change from the DIFFICULTY OpCode to the RANDOM OpCode (RANDOM replaces DIFFICULTY directly), there is no planned change to provide a stronger source of randomness. The RANDOM OpCode relies on the Eth2.0 Randomness Beacon, which has no direct parallel within the context of either JEth or Snowman consensus. Therefore, instead of providing a weaker source of randomness that may be manipulated, the RANDOM OpCode will not be supported. Instead, it will continue the behavior of the DIFFICULTY OpCode of returning the block's difficulty, such that it will always return 1.

## Block Format

### Block Body

- `Version`: provides version of the `ExtData` in the block. Currently, this field is always 0.
- `ExtData`: extra data field within the block body to store atomic transaction bytes.

### Block Header

- `ExtDataHash`: the hash of the bytes in the `ExtDataHash` field
- `BaseFee`: Added by EIP-1559 to represent the base fee of the block (present in Ethereum as of EIP-1559)
- `ExtDataGasUsed`: amount of gas consumed by the atomic transactions in the block
- `BlockGasCost`: surcharge for producing a block faster than the target rate
