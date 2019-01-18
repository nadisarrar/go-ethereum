package vm

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// ExecStats keeps track of accounts and other metrics that are relevant to the execution of a transaction.
type ExecStats struct {
	// set in core/vm/evm.go::NewEVM()
	From          common.Address
	GasPrice      *big.Int
	BlockGasLimit uint64
	Difficulty    *big.Int
	Coinbase      common.Address

	// set in core/state_transition.go::TransitionDb()
	To        common.Address
	Value     *big.Int
	InputLen  uint64
	GasLimit  uint64
	GasUsed   uint64
	GasRefund uint64
	Failed    bool   // true if OOG
	VMError   string // set if Failed == true

	Creates []common.Address // set in core/vm/evm.go::create()
	Writes  []common.Address // set in core/vm/evm.go::Call() and core/vm/instructions.go::opSuicide()
	Reads   []common.Address // set in core/vm/evm.go::CallCode()/DelegateCall()/StaticCall() and in core/vm/instructions.go::opBalance()/opExtCodeSize()/opExtCodeCopy()/opExtCodeHash()/opSuicide()

	// Auxiliary information set in core/state_processor.go::Process()
	Number           *big.Int
	Time             *big.Int
	TxCount          int
	TxIndex          int
	TxHash           common.Hash
	Nonce            uint64
	ReceiptSize      common.StorageSize
	ReceiptStatus    uint64
	ReceiptLogsCount int
}

// Finalize should be run after collecting but before exporting the data
func (execStats *ExecStats) Finalize() {
	execStats.Creates = uniquify(execStats.Creates)
	execStats.Writes = uniquify(execStats.Writes)
	execStats.Reads = uniquify(execStats.Reads)
}

func uniquify(s []common.Address) []common.Address {
	seen := make(map[common.Address]bool, len(s))
	j := 0
	for _, v := range s {
		if _, found := seen[v]; found {
			continue
		}
		seen[v] = true
		s[j] = v
		j++
	}
	return s[:j]
}

func (execStats ExecStats) String() string {
	ret, _ := json.Marshal(execStats)
	return string(ret)
}
