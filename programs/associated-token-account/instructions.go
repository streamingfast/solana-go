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

package associatedtokenaccount

import (
	"bytes"
	"fmt"

	bin "github.com/streamingfast/binary"
	"github.com/streamingfast/solana-go"
	"github.com/streamingfast/solana-go/programs/system"
	"github.com/streamingfast/solana-go/programs/token"
	"github.com/streamingfast/solana-go/text"
)

var PROGRAM_ID = solana.MustPublicKeyFromBase58("ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL")

func init() {
	solana.RegisterInstructionDecoder(PROGRAM_ID, registryDecodeInstruction)
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
		return nil, fmt.Errorf("unable to decode instruction for serum program: %w", err)
	}
	if v, ok := inst.Impl.(solana.AccountSettable); ok {
		err := v.SetAccounts(accounts)
		if err != nil {
			return nil, fmt.Errorf("unable to set accounts for instruction: %w", err)
		}
	}

	return &inst, nil
}

var InstructionDefVariant = bin.NewVariantDefinition(bin.Uint8TypeIDEncoding, []bin.VariantType{
	{"create", (*Create)(nil)},
})

type Instruction struct {
	bin.BaseVariant
}

func (i *Instruction) Accounts() (out []*solana.AccountMeta) {
	switch i.TypeID {
	case 0:
		accounts := i.Impl.(*Create).Accounts
		out = []*solana.AccountMeta{
			accounts.FundingAccount,
			accounts.AssociatedTokenAccount,
			accounts.AssociatedTokenAccountWallet,
			accounts.Mint,
			accounts.SystemProgram,
			accounts.SPLTokenProgram,
			accounts.RentProgram,
		}
	}
	return
}

func (i *Instruction) ProgramID() solana.PublicKey {
	return PROGRAM_ID
}

func (i *Instruction) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := bin.NewEncoder(buf).Encode(i); err != nil {
		return nil, fmt.Errorf("unable to encode instruction: %w", err)
	}
	return buf.Bytes(), nil
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
func (i *Instruction) TextEncode(encoder *text.Encoder, option *text.Option) error {
	return encoder.Encode(i.Impl, option)
}

type CreateAccounts struct {
	FundingAccount               *solana.AccountMeta
	AssociatedTokenAccount       *solana.AccountMeta
	AssociatedTokenAccountWallet *solana.AccountMeta
	Mint                         *solana.AccountMeta
	SystemProgram                *solana.AccountMeta
	SPLTokenProgram              *solana.AccountMeta
	RentProgram                  *solana.AccountMeta
}
type Create struct {
	Accounts *CreateAccounts `bin:"-"`
}

func NewCreateInstruction(
	fundingAccount solana.PublicKey,
	associatedTokenAccount solana.PublicKey,
	associatedTokenAccountWallet solana.PublicKey,
	mint solana.PublicKey,
) *Instruction {
	return &Instruction{
		BaseVariant: bin.BaseVariant{
			TypeID: 0,
			Impl: &Create{
				Accounts: &CreateAccounts{
					FundingAccount:               &solana.AccountMeta{PublicKey: fundingAccount, IsWritable: true, IsSigner: true},
					AssociatedTokenAccount:       &solana.AccountMeta{PublicKey: associatedTokenAccount, IsWritable: true},
					AssociatedTokenAccountWallet: &solana.AccountMeta{PublicKey: associatedTokenAccountWallet},
					Mint:                         &solana.AccountMeta{PublicKey: mint},
					SystemProgram:                &solana.AccountMeta{PublicKey: system.PROGRAM_ID},
					SPLTokenProgram:              &solana.AccountMeta{PublicKey: token.PROGRAM_ID},
					RentProgram:                  &solana.AccountMeta{PublicKey: system.SYSVAR_RENT},
				},
			},
		},
	}
}

func MustGetAssociatedTokenAddress(mint solana.PublicKey, programId solana.PublicKey, owner solana.PublicKey) solana.PublicKey {
	path := [][]byte{
		owner[:],
		programId[:],
		mint[:],
	}

	pubKey, _, err := solana.PublicKeyFindProgramAddress(path, PROGRAM_ID)
	if err != nil {
		panic(fmt.Errorf("unable ot find pda for spl token association: %w", err))
	}
	return pubKey

}
