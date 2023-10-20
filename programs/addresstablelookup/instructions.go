package addresstablelookup

import (
	"bytes"
	"fmt"
	bin "github.com/streamingfast/binary"
	"github.com/streamingfast/solana-go"
)

var ADDRESS_LOOKUP_TABLE_PROGRAM_ID = solana.MustPublicKeyFromBase58("AddressLookupTab1e1111111111111111111111111")

var AddressLookupTableExtendTableInstruction = []byte{02, 00, 00, 00}

func init() {
	solana.RegisterInstructionDecoder(ADDRESS_LOOKUP_TABLE_PROGRAM_ID, registryDecodeInstruction)
}

func registryDecodeInstruction(accounts []*solana.AccountMeta, data []byte) (interface{}, error) {
	inst, err := DecodeInstruction(accounts, data)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

func DecodeInstruction(accounts []*solana.AccountMeta, data []byte) (*Instruction, error) {
	var inst Instruction
	if err := bin.NewDecoder(data).Decode(&inst); err != nil {
		return nil, fmt.Errorf("unable to decode instruction for address lookup table program: %w", err)
	}

	return &inst, nil
}

var InstructionDefVariant = bin.NewVariantDefinition(bin.Uint8TypeIDEncoding, []bin.VariantType{
	{"extendLookupTable", (*ExtendLookupTable)(nil)},
})

type Instruction struct {
	bin.BaseVariant
}

func (i *Instruction) UnmarshalBinary(decoder *bin.Decoder) (err error) {
	return i.BaseVariant.UnmarshalBinaryVariant(decoder, InstructionDefVariant)
}

type ExtendLookupTable struct {
	Index        int32
	AddressCount uint64 `bin:"sizeof=Addresses"`
	Addresses    [][32]byte
}

func ExtendAddressTableLookupInstruction(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	return bytes.Equal(data[:4], AddressLookupTableExtendTableInstruction)
}

func ParseNewAccounts(accounts []byte) [][]byte {
	var newAccounts [][]byte
	numberOfAccounts := len(accounts) / 32

	for i := 0; i < numberOfAccounts; i++ {
		if i == numberOfAccounts {
			break
		}

		addr := accounts[(i * 32) : (i+1)*32]
		newAccounts = append(newAccounts, addr)
	}

	return newAccounts
}
