// Copyright 2020 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpc

import (
	"encoding/json"
	"fmt"
	bin "github.com/streamingfast/binary"
	"github.com/streamingfast/solana-go"
	"github.com/ybbus/jsonrpc"
)

type Context struct {
	Slot bin.Uint64
}

type RPCContext struct {
	Context Context `json:"context,omitempty"`
}

type RPCFilter struct {
	Memcmp   *RPCFilterMemcmp `json:"memcmp,omitempty"`
	DataSize bin.Uint64       `json:"dataSize,omitempty"`
}

type RPCFilterMemcmp struct {
	Offset int           `json:"offset"`
	Bytes  solana.Base58 `json:"bytes"`
}

type SendTransactionOptions struct {
	SkipPreflight       bool           // disable transaction verification step
	PreflightCommitment CommitmentType // preflight commitment level; default: "finalized"
}

// CommitmentType is the level of commitment desired when querying state.
// https://docs.solana.com/developing/clients/jsonrpc-api#configuring-state-commitment
type CommitmentType string

const (
	// CommitmentProcessed queries the most recent block which has reached 1 confirmation by the connected node
	CommitmentProcessed = CommitmentType("processed")
	// CommitmentConfirmed queries the most recent block which has reached 1 confirmation by the cluster
	CommitmentConfirmed = CommitmentType("confirmed")
	// CommitmentConfirmed queries the most recent block which has been finalized by the cluster
	CommitmentFinalized = CommitmentType("finalized")

	// The following are deprecated
	CommitmentMax          = CommitmentType("max")          // Deprecated as of v1.5.5
	CommitmentRecent       = CommitmentType("recent")       // Deprecated as of v1.5.5
	CommitmentRoot         = CommitmentType("root")         // Deprecated as of v1.5.5
	CommitmentSingle       = CommitmentType("single")       // Deprecated as of v1.5.5
	CommitmentSingleGossip = CommitmentType("singleGossip") // Deprecated as of v1.5.5
)

type RpcError struct {
	*jsonrpc.RPCError
	trxError *TransactionError
	Logs     []string
}

func fromRPCError(rerr *jsonrpc.RPCError) *RpcError {
	rpcError := &RpcError{RPCError: rerr}
	v, ok := rpcError.Data.(map[string]interface{})
	if !ok {
		return rpcError
	}
	rpcError.trxError = &TransactionError{Raw: v}
	if err, ok := v["err"].(map[string]interface{}); ok {
		if instructionError, ok := err["InstructionError"].([]interface{}); ok {
			if len(instructionError) == 2 {
				if idx, ok := instructionError[0].(uint64); ok {
					rpcError.trxError.InstructionIndex = idx
				}
				if instErr, ok := instructionError[1].(map[string]interface{}); ok {
					for instErrType, instErrCode := range instErr {
						rpcError.trxError.InstructionErrorType = instErrType

						if str, ok := instErrCode.(string); ok {
							rpcError.trxError.InstructionErrorCode = str
						} else if num, ok := instErrCode.(json.Number); ok {
							rpcError.trxError.InstructionErrorCode = fmt.Sprintf("%s", num)
						} else {
							rpcError.trxError.InstructionErrorCode = "unknown"
						}
					}
				}
			}
		}
	}

	if logs, ok := v["logs"].([]interface{}); ok {
		for _, log := range logs {
			rpcError.Logs = append(rpcError.Logs, log.(string))
		}
	}
	return rpcError
}
