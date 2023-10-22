package addresstablelookup

import (
	"bytes"
	"fmt"
	"github.com/mr-tron/base58"
	bin "github.com/streamingfast/binary"
	"github.com/streamingfast/solana-go"
	"strings"
)

var AddressLookupTableProgramId = solana.MustPublicKeyFromBase58("AddressLookupTab1e1111111111111111111111111")

func init() {
	solana.RegisterInstructionDecoder(AddressLookupTableProgramId, registryDecodeInstruction)
}

func registryDecodeInstruction(_ []*solana.AccountMeta, data []byte) (interface{}, error) {
	inst, err := DecodeInstruction(data)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

func DecodeInstruction(data []byte) (*Instruction, error) {
	var inst Instruction
	if err := bin.NewDecoder(data).Decode(&inst); err != nil {
		return nil, fmt.Errorf("unable to decode instruction for address lookup table program: %w", err)
	}

	return &inst, nil
}

var InstructionDefVariant = bin.NewVariantDefinition(bin.Uint32TypeIDEncoding, []bin.VariantType{
	{"createTableLookup", (*CreateLookupTable)(nil)},
	{"freezeLookupTable", (*FreezeLookupTable)(nil)},
	{"extendLookupTable", (*ExtendLookupTable)(nil)},
	{"deactivateLookupTable", (*DeactivateLookupTable)(nil)},
	{"closeLookupTable", (*CloseLookupTable)(nil)},
})

type Instruction struct {
	bin.BaseVariant
}

func (i *Instruction) UnmarshalBinary(decoder *bin.Decoder) (err error) {
	return i.BaseVariant.UnmarshalBinaryVariant(decoder, InstructionDefVariant)
}

func (i *Instruction) MarshalBinary(encoder *bin.Encoder) error {
	err := encoder.WriteUint8(uint8(i.TypeID))
	if err != nil {
		return fmt.Errorf("unable to write variant type: %w", err)
	}
	return encoder.Encode(i.Impl)
}

type CreateLookupTable struct {
	RecentSlot uint64
	BumpSeed   uint8
}

type FreezeLookupTable struct{}

type ExtendLookupTable struct {
	AddressCount uint64 `bin:"sizeof=Addresses"`
	Addresses    [][32]byte
}

func (e *ExtendLookupTable) String() string {
	sb := strings.Builder{}
	sb.WriteString("{\n")
	for _, addr := range e.Addresses {
		sb.WriteString("\t")
		sb.WriteString(base58.Encode(addr[:]))
		sb.WriteString("\n")
	}
	sb.WriteString("}\n")
	return sb.String()
}

type DeactivateLookupTable struct{}

type CloseLookupTable struct{}

func (i *Instruction) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := bin.NewEncoder(buf).Encode(i); err != nil {
		return nil, fmt.Errorf("unable to encode instruction: %w", err)
	}
	return buf.Bytes(), nil
}
