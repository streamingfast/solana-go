package solana

import (
	"crypto"
	"crypto/ed25519"
	crypto_rand "crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"filippo.io/edwards25519"
	"github.com/mr-tron/base58"
)

const MAX_SEED_LENGTH = 32

type PrivateKey []byte

func MustPrivateKeyFromBase58(in string) PrivateKey {
	out, err := PrivateKeyFromBase58(in)
	if err != nil {
		panic(err)
	}
	return out
}

func PrivateKeyFromBase58(privkey string) (PrivateKey, error) {
	res, err := base58.Decode(privkey)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func PrivateKeyFromSolanaKeygenFile(file string) (PrivateKey, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read keygen file: %w", err)
	}

	var values []uint8
	err = json.Unmarshal(content, &values)
	if err != nil {
		return nil, fmt.Errorf("decode keygen file: %w", err)
	}

	return values, nil
}

func (k PrivateKey) String() string {
	return base58.Encode(k)
}

func NewRandomPrivateKey() (PublicKey, PrivateKey, error) {
	pub, priv, err := ed25519.GenerateKey(crypto_rand.Reader)
	if err != nil {
		return PublicKey{}, nil, err
	}
	var publicKey PublicKey
	copy(publicKey[:], pub)
	return publicKey, PrivateKey(priv), nil
}

func (k PrivateKey) Sign(payload []byte) (Signature, error) {
	p := ed25519.PrivateKey(k)
	signData, err := p.Sign(crypto_rand.Reader, payload, crypto.Hash(0))
	if err != nil {
		return Signature{}, err
	}

	var signature Signature
	copy(signature[:], signData)

	return signature, err
}

func (k PrivateKey) PublicKey() PublicKey {
	p := ed25519.PrivateKey(k)
	pub := p.Public().(ed25519.PublicKey)

	var publicKey PublicKey
	copy(publicKey[:], pub)

	return publicKey
}

type PublicKey [32]byte

func PublicKeyFromBytes(in []byte) (out PublicKey) {
	byteCount := len(in)
	if byteCount == 0 {
		return
	}

	max := 32
	if byteCount < max {
		max = byteCount
	}

	copy(out[:], in[0:max])
	return
}

func MustPublicKeyFromBase58(in string) PublicKey {
	out, err := PublicKeyFromBase58(in)
	if err != nil {
		panic(err)
	}
	return out
}

func PublicKeyFromBase58(in string) (out PublicKey, err error) {
	val, err := base58.Decode(in)
	if err != nil {
		return out, fmt.Errorf("decode: %w", err)
	}

	if len(val) != 32 {
		return out, fmt.Errorf("invalid length, expected 32, got %d", len(val))
	}

	copy(out[:], val)
	return
}

func PublicKeyFindProgramAddress(path [][]byte, programId PublicKey) (PublicKey, byte, error) {
	nonce := byte(255)
	for {
		seedsWithNonce := append(path, []byte{nonce})
		key, err := createProgramAddress(seedsWithNonce, programId)
		if err == nil {
			return key, nonce, nil
		}

		nonce -= 1
		if nonce == 0 {
			return PublicKey{}, 0x00, fmt.Errorf("unable to find a viable program address nonce")
		}
	}
}

func createProgramAddress(seeds [][]byte, programId PublicKey) (PublicKey, error) {
	buf := []byte{}
	for _, seed := range seeds {
		if len(seed) > MAX_SEED_LENGTH {
			return PublicKey{}, fmt.Errorf("max seed length exceeded")
		}
		buf = append(buf, seed...)
	}
	buf = append(buf, programId[:]...)
	buf = append(buf, []byte("ProgramDerivedAddress")...)
	pkey := sha256.Sum256(buf)
	if _, err := new(edwards25519.Point).SetBytes(pkey[:]); err == nil {
		return PublicKey{}, fmt.Errorf("invalid seeds, address must fall off the curve")
	}
	return pkey, nil
}
func (p PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(base58.Encode(p[:]))
}

func (p *PublicKey) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*p, err = PublicKeyFromBase58(s)
	if err != nil {
		return fmt.Errorf("invalid public key %q: %w", s, err)
	}
	return
}

func (p PublicKey) Equals(pb PublicKey) bool {
	return p == pb
}

var zeroPublicKey = PublicKey{}

func (p PublicKey) IsZero() bool {
	return p == zeroPublicKey
}

func (p PublicKey) String() string {
	return base58.Encode(p[:])
}
