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
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"reflect"
	"time"

	bin "github.com/streamingfast/binary"
	"github.com/streamingfast/logging"
	"github.com/streamingfast/solana-go"
	"github.com/ybbus/jsonrpc"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ErrNotFound = errors.New("not found")

type Client struct {
	rpcURL             string
	rpcClient          jsonrpc.RPCClient
	headers            http.Header
	requestIDGenerator func() int
}

func NewClient(rpcURL string) *Client {
	return &Client{
		rpcURL: rpcURL,
		rpcClient: jsonrpc.NewClientWithOpts(rpcURL, &jsonrpc.RPCClientOpts{
			HTTPClient: &http.Client{
				Transport: &withLoggingRoundTripper{
					defaultLogger: &zlog,
					tracer:        tracer,
				}},
		}),
		requestIDGenerator: generateRequestID,
	}
}

func (c *Client) SetHeader(k, v string) {
	if c.headers == nil {
		c.headers = http.Header{}
	}
	c.headers.Set(k, v)
}

func (c *Client) GetBalance(ctx context.Context, publicKey string, commitment CommitmentType) (out *GetBalanceResult, err error) {
	commit := map[string]string{
		"commitment": string(commitment),
	}
	params := []interface{}{publicKey}
	if commitment != "" {
		params = append(params, commit)
	}

	err = c.callFor(&out, "getBalance", params...)
	return
}

func (c *Client) GetRecentBlockhash(ctx context.Context, commitment CommitmentType) (out *GetRecentBlockhashResult, err error) {
	commit := map[string]string{
		"commitment": string(commitment),
	}
	var params []interface{}
	if commitment != "" {
		params = append(params, commit)
	}

	err = c.callFor(&out, "getRecentBlockhash", params)
	return
}

func (c *Client) GetSlot(ctx context.Context, commitment CommitmentType) (out GetSlotResult, err error) {
	commit := map[string]string{
		"commitment": string(commitment),
	}
	var params []interface{}
	if commitment != "" {
		params = append(params, commit)
	}

	err = c.callFor(&out, "getSlot", params)
	return
}

func (c *Client) GetConfirmedBlock(ctx context.Context, slot uint64, encoding string) (out *GetConfirmedBlockResult, err error) {
	if encoding == "" {
		encoding = "json"
	}
	params := []interface{}{slot, encoding}

	err = c.callFor(&out, "getConfirmedBlock", params...)
	return
}

func (c *Client) GetAccountInfo(ctx context.Context, account solana.PublicKey) (out *GetAccountInfoResult, err error) {
	obj := map[string]interface{}{
		"encoding": "base64",
	}
	params := []interface{}{account, obj}

	err = c.callFor(&out, "getAccountInfo", params...)
	if err != nil {
		return nil, err
	}

	if out.Value == nil {
		return nil, ErrNotFound
	}

	return out, nil
}

func (c *Client) GetAccountDataIn(ctx context.Context, account solana.PublicKey, inVar interface{}) (err error) {
	resp, err := c.GetAccountInfo(ctx, account)
	if err != nil {
		return err
	}

	return bin.NewDecoder(resp.Value.Data).Decode(inVar)
}

func (c *Client) GetConfirmedTransaction(ctx context.Context, signature string) (out TransactionWithMeta, err error) {
	params := []interface{}{signature, "json"}

	err = c.callFor(&out, "getConfirmedTransaction", params...)
	return
}

func (c *Client) GetConfirmedSignaturesForAddress2(ctx context.Context, address solana.PublicKey, opts *GetConfirmedSignaturesForAddress2Opts) (out GetConfirmedSignaturesForAddress2Result, err error) {

	params := []interface{}{address.String(), opts}

	err = c.callFor(&out, "getConfirmedSignaturesForAddress2", params...)
	return
}

func (c *Client) GetSignaturesForAddress(ctx context.Context, address solana.PublicKey, opts *GetSignaturesForAddressOpts) (out GetSignaturesForAddressResult, err error) {
	params := []interface{}{address.String(), opts}

	err = c.callFor(&out, "getSignaturesForAddress", params...)
	return
}

func (c *Client) GetProgramAccounts(ctx context.Context, publicKey solana.PublicKey, opts *GetProgramAccountsOpts) (out GetProgramAccountsResult, err error) {
	obj := map[string]interface{}{
		"encoding": "base64",
	}
	if opts != nil {
		if opts.Commitment != "" {
			obj["commitment"] = string(opts.Commitment)
		}
		if len(opts.Filters) != 0 {
			obj["filters"] = opts.Filters
		}
	}

	params := []interface{}{publicKey, obj}

	err = c.callFor(&out, "getProgramAccounts", params...)
	return
}

func (c *Client) GetMinimumBalanceForRentExemption(ctx context.Context, dataSize int) (lamport int, err error) {
	params := []interface{}{dataSize}
	err = c.callFor(&lamport, "getMinimumBalanceForRentExemption", params...)
	return
}

type SimulateTransactionResponse struct {
	Err  interface{}
	Logs []string
}

func (c *Client) SimulateTransaction(ctx context.Context, transaction *solana.Transaction) (*SimulateTransactionResponse, error) {
	buf := new(bytes.Buffer)
	if err := bin.NewEncoder(buf).Encode(transaction); err != nil {
		return nil, fmt.Errorf("send transaction: encode transaction: %w", err)
	}
	trxData := buf.Bytes()

	obj := map[string]interface{}{
		"encoding": "base64",
	}

	b64Data := base64.StdEncoding.EncodeToString(trxData)
	params := []interface{}{
		b64Data,
		obj,
	}

	var out *SimulateTransactionResponse
	if err := c.callFor(&out, "simulateTransaction", params...); err != nil {
		return nil, fmt.Errorf("send transaction: rpc send: %w", err)
	}

	return out, nil

}

func (c *Client) SendTransaction(
	ctx context.Context,
	transaction *solana.Transaction,
	opts *SendTransactionOptions,
) (signature string, err error) {
	buf := new(bytes.Buffer)

	if err := bin.NewEncoder(buf).Encode(transaction); err != nil {
		return "", fmt.Errorf("send transaction: encode transaction: %w", err)
	}

	trxData := buf.Bytes()
	fmt.Println("Encodeded trx: ", hex.EncodeToString(trxData))

	obj := map[string]interface{}{
		"encoding": "base64",
	}
	if opts != nil {
		if opts.SkipPreflight {
			obj["skipPreflight"] = opts.SkipPreflight
		}
		if opts.PreflightCommitment != "" {
			obj["preflightCommitment"] = opts.PreflightCommitment
		}
	}

	params := []interface{}{
		base64.StdEncoding.EncodeToString(trxData),
		obj,
	}

	if err := c.callFor(&signature, "sendTransaction", params...); err != nil {
		return "", fmt.Errorf("send transaction: rpc send: %w", err)
	}
	return
}

func (c *Client) RequestAirdrop(ctx context.Context, account *solana.PublicKey, lamport uint64, commitment CommitmentType) (signature string, err error) {

	obj := map[string]interface{}{
		"commitment": commitment,
	}

	params := []interface{}{
		account.String(),
		lamport,
		obj,
	}

	if err := c.callFor(&signature, "requestAirdrop", params...); err != nil {
		return "", fmt.Errorf("send transaction: rpc send: %w", err)
	}
	return
}

func (c *Client) callFor(out interface{}, method string, params ...interface{}) error {
	request := jsonrpc.NewRequest(method, params...)
	request.ID = c.requestIDGenerator()

	logger := zlog.With(zap.Int("id", request.ID), zap.String("method", method))
	ctx := logging.WithLogger(context.Background(), logger)

	fields := []zapcore.Field{}
	if tracer.Enabled() {
		fields = append(fields, zap.Reflect("params", params))
	}
	fields = append(fields, zapType("output", out))

	startTime := time.Now()
	decodingTime := time.Time{}

	logger.Info("performing JSON-RPC call", fields...)
	defer func() {
		fields := []zapcore.Field{}
		if !decodingTime.IsZero() {
			fields = append(fields, zap.Duration("parsing", time.Since(decodingTime)))
		}
		fields = append(fields, zap.Duration("overall", time.Since(startTime)))

		logger.Info("performed JSON-RPC call", fields...)
	}()

	// When `jsonrpc` library we use accepts `ctx contxt.Context` as first parameter in `CallCtxRaw` (or other name),
	// replace with appropriate function that accept context value.
	//
	// See https://github.com/ybbus/jsonrpc/pull/39
	_ = ctx
	rpcResponse, err := c.rpcClient.CallRaw(request)
	if err != nil {
		return fmt.Errorf("call raw: %w", err)
	}

	if rpcResponse.Error != nil {
		fmt.Println("GGGGRRRR:", rpcResponse.Result)
		return fmt.Errorf("rpc response: %w", rpcResponse.Error)
	}

	return rpcResponse.GetObject(out)
}

var requestCounter = atomic.NewInt64(0)

func generateRequestID() int {
	return int(requestCounter.Inc())
}

func zapType(key string, v interface{}) zap.Field {
	return zap.Stringer(key, zapTypeWrapper{v})
}

type zapTypeWrapper struct {
	v interface{}
}

func (w zapTypeWrapper) String() string {
	return reflect.TypeOf(w.v).String()
}

type withLoggingRoundTripper struct {
	defaultLogger **zap.Logger
	tracer        logging.Tracer
}

func (t *withLoggingRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	logger := logging.Logger(request.Context(), *t.defaultLogger)

	debugEnabled := logger.Core().Enabled(zap.DebugLevel)
	traceEnabled := t.tracer.Enabled()

	if debugEnabled {
		requestDump, err := httputil.DumpRequestOut(request, true)
		if err != nil {
			panic(fmt.Errorf("unexpecting that httputil.DumpRequestOut would panic: %w", err))
		}

		logger.Debug("JSON-RPC request\n" + string(requestDump))
	}

	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	if debugEnabled {
		responseDump, err := httputil.DumpResponse(response, traceEnabled)
		if err != nil {
			panic(fmt.Errorf("unexpecting that httputil.DumpRequestOut would panic: %w", err))
		}

		logger.Debug("JSON-RPC response\n" + string(responseDump))
	}

	return response, nil
}
