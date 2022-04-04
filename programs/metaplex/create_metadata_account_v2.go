package metaplex

import (
	"github.com/streamingfast/solana-go"
	"github.com/streamingfast/solana-go/programs/system"
)

type CreateMetadataV2Account struct {
	Instruction InstType
	Data        DataV2
	IsMutable   bool

	Accounts *CreateMetadataAccountV2Accounts `borsh_skip:"true"`
}

func (i CreateMetadataV2Account) ListAccounts() []*solana.AccountMeta {
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

/// Create Metadata object.
///   0. `[writable]`  Metadata key (pda of ['metadata', program id, mint id])
///   1. `[]` Mint of token asset
///   2. `[signer]` Mint authority
///   3. `[signer]` payer
///   4. `[]` update authority info
///   5. `[]` System program
///   6. `[]` Rent info
type CreateMetadataAccountV2Accounts struct {
	Metadata        *solana.AccountMeta
	Mint            *solana.AccountMeta
	MintAuthority   *solana.AccountMeta
	Payer           *solana.AccountMeta
	UpdateAuthority *solana.AccountMeta
	SystemProgram   *solana.AccountMeta
	RentProgram     *solana.AccountMeta
}

func NewCreateMetadataAccountV2Instruction(
	programID solana.PublicKey,
	data DataV2,
	isMutable bool,
	metadata,
	mint,
	mintAuthority,
	payer,
	updateAuthority solana.PublicKey,
) *Instruction {
	var inst = CreateMetadataV2Account{
		Instruction: CreateMetadataAccountV2Inst,
		Data:        data,
		IsMutable:   isMutable,
		Accounts: &CreateMetadataAccountV2Accounts{
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
