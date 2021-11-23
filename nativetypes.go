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

package solana

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/mr-tron/base58"
)

type Padding []byte

type Hash PublicKey

type Signature [64]byte

// Deprecated: Use NewSignatureFromBase58 instead
func SignatureFromBase58(in string) (out Signature, err error) {
	return NewSignatureFromBase58(in)
}

func NewSignatureFromBytes(in []byte) (out Signature, err error) {
	if len(in) != 64 {
		err = fmt.Errorf("invalid length, expected 64, got %d", len(in))
		return
	}
	copy(out[:], in)
	return
}

func NewSignatureFromString(in string) (out Signature, err error) {
	bytes, err := hex.DecodeString(in)
	if err != nil {
		return out, fmt.Errorf("hex decode: %w", err)
	}

	return NewSignatureFromBytes(bytes)
}

func MustSignatureFromString(in string) (out Signature) {
	out, err := NewSignatureFromString(in)
	if err != nil {
		panic(err)
	}
	return
}

func NewSignatureFromBase58(in string) (out Signature, err error) {
	bytes, err := base58.Decode(in)
	if err != nil {
		return out, fmt.Errorf("base58 decode: %w", err)
	}

	return NewSignatureFromBytes(bytes)
}
func (s Signature) ToSlice() []byte {
	return s[:]
}

func (s Signature) MarshalJSON() ([]byte, error) {
	return json.Marshal(base58.Encode(s[:]))
}

func (s *Signature) UnmarshalJSON(data []byte) (err error) {
	var str string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	dat, err := base58.Decode(str)
	if err != nil {
		return err
	}

	if len(dat) != 64 {
		return errors.New("invalid data length for public key")
	}

	target := Signature{}
	copy(target[:], dat)
	*s = target
	return
}

func (s Signature) Verify(publicKey PublicKey, message []byte) bool {
	return ed25519.Verify(ed25519.PublicKey(publicKey[:]), message, s[:])
}

func (s Signature) String() string {
	return base58.Encode(s[:])
}

///
type Base58 []byte

func (t Base58) MarshalJSON() ([]byte, error) {
	return json.Marshal(base58.Encode(t))
}

func (t *Base58) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	*t, err = base58.Decode(s)
	return
}

func (t Base58) String() string {
	return base58.Encode(t)
}

type Data []byte

func (t Data) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"data":     []byte(t),
		"encoding": "base64",
	})
}

func (t *Data) UnmarshalJSON(data []byte) (err error) {
	var in []string
	if err := json.Unmarshal(data, &in); err != nil {
		return err
	}

	if len(in) != 2 {
		return fmt.Errorf("invalid length for solana.Data, expected 2, found %d", len(in))
	}

	switch in[1] {
	case "base64":
		*t, err = base64.StdEncoding.DecodeString(in[0])
	default:
		return fmt.Errorf("unsupported encoding %s", in[1])
	}
	return
}

func (t Data) String() string {
	return base64.StdEncoding.EncodeToString(t)
}

///
type ByteWrapper struct {
	io.Reader
}

func (w *ByteWrapper) ReadByte() (byte, error) {
	var b [1]byte
	_, err := w.Read(b[:])
	return b[0], err
}
