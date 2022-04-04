package metaplex

import (
	"github.com/streamingfast/solana-go"
	"github.com/streamingfast/solana-go/programs/system"
	"github.com/streamingfast/solana-go/programs/token"
)

type CreateMasterEditionV3Account struct {
	Instruction InstType
	MaxSupply   *uint64

	Accounts *CreateMetadataMasterEditionV3Accounts `borsh_skip:"true"`
}

func (i CreateMasterEditionV3Account) ListAccounts() []*solana.AccountMeta {
	return []*solana.AccountMeta{
		i.Accounts.Edition,
		i.Accounts.Mint,
		i.Accounts.UpdateAuthority,
		i.Accounts.MintAuthority,
		i.Accounts.Payer,
		i.Accounts.Metadata,
		i.Accounts.TokenProgram,
		i.Accounts.SystemProgram,
		i.Accounts.RentProgram,
	}
}

/// Register a Metadata as a Master Edition V2, which means Edition V2s can be minted.
/// Henceforth, no further tokens will be mintable from this primary mint. Will throw an error if more than one
/// token exists, and will throw an error if less than one token exists in this primary mint.
///   0. `[writable]` Unallocated edition V2 account with address as pda of ['metadata', program id, mint, 'edition']
///   1. `[writable]` Metadata mint
///   2. `[signer]` Update authority
///   3. `[signer]` Mint authority on the metadata's mint - THIS WILL TRANSFER AUTHORITY AWAY FROM THIS KEY
///   4. `[signer]` payer
///   5.  [writable] Metadata account
///   6. `[]` Token program
///   7. `[]` System program
///   8. `[]` Rent info
type CreateMetadataMasterEditionV3Accounts struct {
	Edition         *solana.AccountMeta
	Mint            *solana.AccountMeta
	UpdateAuthority *solana.AccountMeta
	MintAuthority   *solana.AccountMeta
	Payer           *solana.AccountMeta
	Metadata        *solana.AccountMeta
	TokenProgram    *solana.AccountMeta
	SystemProgram   *solana.AccountMeta
	RentProgram     *solana.AccountMeta
}

func NewCreateMetadataMasterEditionV3Instruction(
	programID solana.PublicKey,
	maxSupply *uint64,
	edition,
	mint,
	updateAuthority,
	mintAuthority,
	payer,
	metadata solana.PublicKey,
) *Instruction {
	var inst = CreateMasterEditionV3Account{
		Instruction: CreateMasterEditionV3Inst,
		MaxSupply:   maxSupply,
		Accounts: &CreateMetadataMasterEditionV3Accounts{

			Edition:         &solana.AccountMeta{PublicKey: edition, IsWritable: true},
			Mint:            &solana.AccountMeta{PublicKey: mint, IsWritable: true},
			UpdateAuthority: &solana.AccountMeta{PublicKey: updateAuthority, IsSigner: true},
			MintAuthority:   &solana.AccountMeta{PublicKey: mintAuthority, IsSigner: true},
			Payer:           &solana.AccountMeta{PublicKey: payer, IsSigner: true},
			Metadata:        &solana.AccountMeta{PublicKey: metadata, IsWritable: true},
			TokenProgram:    &solana.AccountMeta{PublicKey: token.PROGRAM_ID},
			SystemProgram:   &solana.AccountMeta{PublicKey: system.PROGRAM_ID},
			RentProgram:     &solana.AccountMeta{PublicKey: system.SYSVAR_RENT},
		},
	}
	return NewInstruction(programID, inst)
}
