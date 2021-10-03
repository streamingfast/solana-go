package metaplex

import (
	"github.com/streamingfast/solana-go"
	"github.com/streamingfast/solana-go/programs/system"
)

type CreateMeta struct {
	Instruction InstType
	Data        Data
	IsMutable   bool

	Accounts *CreateMetaAccounts `borsh_skip:"true"`
}

func (i CreateMeta) ListAccounts() []*solana.AccountMeta {
	return []*solana.AccountMeta{
		i.Accounts.Metadata,
		i.Accounts.Mint,
		i.Accounts.MintAuthority,
		i.Accounts.Payer,
		i.Accounts.UpdateAuthority,
		i.Accounts.SystemProgram,
		i.Accounts.RentProgram,
	}
}

type CreateMetaAccounts struct {
	Metadata        *solana.AccountMeta
	Mint            *solana.AccountMeta
	MintAuthority   *solana.AccountMeta
	Payer           *solana.AccountMeta
	UpdateAuthority *solana.AccountMeta
	SystemProgram   *solana.AccountMeta
	RentProgram     *solana.AccountMeta
}

// Create Metadata object.
//   0. `[writable]`  Metadata key (pda of ['metadata', program id, mint id])
//   1. `[]` Mint of token asset
//   2. `[signer]` Mint authority
//   3. `[signer]` payer
//   4. `[]` update authority info
//   5. `[]` System program
//   6. `[]` Rent info

func NewCreateMetaInstruction(programID solana.PublicKey, data Data, isMutable bool, metadata, mint, mintAuthority, payer, updateAuthority solana.PublicKey) *Instruction {
	var inst = CreateMeta{
		Instruction: CreateMetadataObjectIns,
		Data:        data,
		IsMutable:   isMutable,
		Accounts: &CreateMetaAccounts{
			Metadata:        &solana.AccountMeta{PublicKey: metadata, IsWritable: true},
			Mint:            &solana.AccountMeta{PublicKey: mint},
			MintAuthority:   &solana.AccountMeta{PublicKey: mintAuthority, IsSigner: true},
			Payer:           &solana.AccountMeta{PublicKey: payer, IsSigner: true},
			UpdateAuthority: &solana.AccountMeta{PublicKey: updateAuthority},
			SystemProgram:   &solana.AccountMeta{PublicKey: system.PROGRAM_ID},
			RentProgram:     &solana.AccountMeta{PublicKey: system.SYSVAR_RENT},
		},
	}
	return NewInstruction(programID, inst)
}
