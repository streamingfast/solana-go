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

type ClientOption = func(cli *Client) *Client

var WithDebug = func() ClientOption {
	return func(cli *Client) *Client {
		cli.debug = true
		return cli
	}
}

type Client struct {
	rpcURL             string
	rpcClient          jsonrpc.RPCClient
	headers            http.Header
	requestIDGenerator func() int
	debug              bool
}

func NewClient(rpcURL string, opts ...ClientOption) *Client {
	c := &Client{
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

	for _, opt := range opts {
		c = opt(c)
	}

	return c
}

func (c *Client) SetHeader(k, v string) {
	if c.headers == nil {
		c.headers = http.Header{}
	}
	c.headers.Set(k, v)
}

func (c *Client) SendTransaction(transaction *solana.Transaction, opts *SendTransactionOptions) (signature string, err error) {
	buf := new(bytes.Buffer)

	if err := bin.NewEncoder(buf).Encode(transaction); err != nil {
		return "", fmt.Errorf("send transaction: encode transaction: %w", err)
	}

	trxData := buf.Bytes()

	obj := map[string]interface{}{
		"encoding": "base64",
	}
	if opts != nil {
		if opts.SkipPreflight {
			obj["skipPreflight"] = opts.SkipPreflight
		}
		if opts.PreflightCommitment != "" {
			obj["preflightCaommitment"] = opts.PreflightCommitment
		}
	}

	params := []interface{}{
		base64.StdEncoding.EncodeToString(trxData),
		obj,
	}

	if err := c.DoRequest(&signature, "sendTransaction", params...); err != nil {
		var rpcError *jsonrpc.RPCError
		if errors.As(err, &rpcError) {
			instructionError := fromRPCError(rpcError)
			if c.debug && instructionError.trxError != nil {
				fmt.Println("RPC ERROR")
				fmt.Printf("Instruction Index %d error: %s -> %s\n", instructionError.trxError.InstructionIndex, instructionError.trxError.InstructionErrorType, instructionError.trxError.InstructionErrorCode)
				for _, log := range instructionError.Logs {
					fmt.Println("> ", log)
				}
				zlog.Info("encountered RPC error", zap.Reflect("instruction_error", instructionError))
			}
		}
		return "", fmt.Errorf("send transaction: rpc send: %w", err)
	}
	return
}

func (c *Client) DoRequest(out interface{}, method string, params ...interface{}) error {
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

	logger.Debug("performing JSON-RPC call", fields...)
	defer func() {
		fields := []zapcore.Field{}
		if !decodingTime.IsZero() {
			fields = append(fields, zap.Duration("parsing", time.Since(decodingTime)))
		}
		fields = append(fields, zap.Duration("overall", time.Since(startTime)))

		logger.Debug("performed JSON-RPC call", fields...)
	}()

	_ = ctx
	rpcResponse, err := c.rpcClient.CallRaw(request)
	if err != nil {
		return fmt.Errorf("call raw: %w", err)
	}

	if rpcResponse.Error != nil {
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

	if debugEnabled {
		if tracer.Enabled() {
			requestDump, err := httputil.DumpRequestOut(request, true)
			if err != nil {
				panic(fmt.Errorf("unexpecting that httputil.DumpRequestOut would panic: %w", err))
			}

			logger.Debug("JSON-RPC request\n" + string(requestDump))
		} else {
			logger.Debug(fmt.Sprintf("JSON-RPC request %s %s", request.Method, request.URL.String()))
		}
	}

	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	if debugEnabled {
		if tracer.Enabled() {
			responseDump, err := httputil.DumpResponse(response, true)
			if err != nil {
				panic(fmt.Errorf("unexpecting that httputil.DumpRequestOut would panic: %w", err))
			}

			logger.Debug("JSON-RPC response\n" + string(responseDump))
		} else {
			logger.Debug(fmt.Sprintf("JSON-RPC response %s (%d bytes)", response.Status, response.ContentLength))
		}
	}

	return response, nil
}
