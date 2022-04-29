package rpc

import (
	"encoding/json"
	"fmt"
	bin "github.com/streamingfast/binary"
	"github.com/streamingfast/solana-go"
)

type GetTransactionResponse struct {
	Slot        bin.Uint64   `json:"slot"`
	BlockTime   *bin.Uint64  `json:"blockTime"`
	Transaction *Transaction `json:"transaction"`
	Meta        *Meta        `json:"meta"`
}

type Transaction struct {
	Signatures []string `json:"signatures"`
	Message    *Message `json:"message"`
}

type Message struct {
	AccountKeys     []solana.PublicKey `json:"accountKeys"`
	Header          MessageHeader      `json:"header"`
	Instructions    []Instruction      `json:"instructions"`
	RecentBlockhash solana.PublicKey   `json:"recentBlockhash"`
}

type MessageHeader struct {
	NumReadonlySignedAccounts   bin.Uint64 `json:"numReadonlySignedAccounts"`
	NumReadonlyUnsignedAccounts bin.Uint64 `json:"numReadonlyUnsignedAccounts"`
	NumRequiredSignatures       bin.Uint64 `json:"numRequiredSignatures"`
}

type Instruction struct {
	ProgramIdIndex bin.Uint64   `json:"programIdIndex"`
	Accounts       []bin.Uint64 `json:"accounts"`
	Data           string       `json:"data"`
}

type Meta struct {
	Err               *TransactionError   `json:"err"`
	Fee               bin.Uint64          `json:"fee"`
	PreBalances       []bin.Uint64        `json:"preBalances"`
	PostBalances      []bin.Uint64        `json:"postBalances"`
	InnerInstructions []*InnerInstruction `json:"innerInstructions"`
	PostTokenBalances []*TokeBalance      `json:"postTokenBalances"`
	PreTokenBalances  []*TokeBalance      `json:"preTokenBalances"`
	LogMessages       []string            `json:"logMessages"`
	Rewards           []interface{}       `json:"rewards"`
}

type InnerInstruction struct {
	Index        bin.Uint64        `json:"index"`
	Instructions []InstructionMeta `json:"instructions"`
}

type InstructionMeta struct {
	Accounts       []bin.Uint64 `json:"accounts"`
	Data           string       `json:"data"`
	ProgramIdIndex bin.Uint64   `json:"programIdIndex"`
}

type TokeBalance struct {
	AccountIndex bin.Uint64       `json:"accountIndex"`
	Mint         solana.PublicKey `json:"mint"`
	Owner        solana.PublicKey `json:"owner"`
	//UiTokenAmount struct {
	//	Amount   string     `json:"amount"`
	//	Decimals bin.Uint64 `json:"decimals"`
	//	// deprecated
	//	UiAmount       *float64 `json:"uiAmount"`
	//	UiAmountString string   `json:"uiAmountString"`
	//} `json:"uiTokenAmount"`
}

type TransactionError struct {
	Raw                  map[string]interface{} `json:"data,omitempty"`
	InstructionIndex     uint64
	InstructionErrorCode string
	InstructionErrorType string
}

func (t *TransactionError) UnmarshalJSON(data []byte) (err error) {
	var errMap map[string]interface{}
	if err := json.Unmarshal(data, &errMap); err != nil {
		return err
	}
	t.Raw = errMap
	if instructionError, ok := t.Raw["InstructionError"].([]interface{}); ok {
		if len(instructionError) == 2 {
			if idx, ok := instructionError[0].(uint64); ok {
				t.InstructionIndex = idx
			}
			if instErr, ok := instructionError[1].(map[string]interface{}); ok {
				for instErrType, instErrCode := range instErr {
					t.InstructionErrorType = instErrType

					if str, ok := instErrCode.(string); ok {
						t.InstructionErrorCode = str
					} else if num, ok := instErrCode.(json.Number); ok {
						t.InstructionErrorCode = fmt.Sprintf("%s", num)
					} else {
						t.InstructionErrorCode = "unknown"
					}
				}
			}
		}
	}
	return
}

func (c *Client) GetConfirmedTransaction(signature string) (out *GetTransactionResponse, err error) {
	conf := CommitmentConfirmed
	return c.GetTransaction(signature, &conf)
}

// GetTransaction For processing many dependent transactions in series, it's recommended to use "confirmed" commitment, which balances speed with rollback safety. For total safety, it's recommended to use"finalized" commitment.
func (c *Client) GetTransaction(signature string, commitmentType *CommitmentType) (out *GetTransactionResponse, err error) {
	opts := map[string]interface{}{
		"encoding": "json",
	}
	if commitmentType != nil {
		opts["Commitment"] = *commitmentType
	}
	params := []interface{}{signature, opts}
	err = c.DoRequest(&out, "getTransaction", params...)
	return
}
