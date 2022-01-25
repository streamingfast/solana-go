package metaplex

import (
	"github.com/streamingfast/solana-go"
)

type UpdateMetadataAccount struct {
	Instruction         InstType
	Data                *Data
	UpdateAuthority     *solana.PublicKey
	PrimarySaleHappened *bool

	Accounts *UpdateMetadataAccountAccounts `borsh_skip:"true"`
}

func (i UpdateMetadataAccount) ListAccounts() []*solana.AccountMeta {
	return []*solana.AccountMeta{
		i.Accounts.Metadata,
		i.Accounts.UpdateAuthorityKey,
	}
}

/// Update a Metadata
///   0. `[writable]` Metadata account
///   1. `[signer]` Update authority key
type UpdateMetadataAccountAccounts struct {
	Metadata           *solana.AccountMeta
	UpdateAuthorityKey *solana.AccountMeta
}

// Create Metadata object.
//   0. `[writable]`  Metadata key (pda of ['metadata', program id, mint id])
//   1. `[]` Mint of token asset
//   2. `[signer]` Mint authority
//   3. `[signer]` payer
//   4. `[]` update authority info
//   5. `[]` System program
//   6. `[]` Rent info

func NewUpdateMetadataAccountInstruction(
	programID solana.PublicKey,
	data *Data,
	updateAuthority *solana.PublicKey,
	primarySaleHappened *bool,
	metadata solana.PublicKey,
	updateAuthorityKey solana.PublicKey,
) *Instruction {
	var inst = UpdateMetadataAccount{
		Instruction:         UpdateMetadataAccountInst,
		Data:                data,
		UpdateAuthority:     updateAuthority,
		PrimarySaleHappened: primarySaleHappened,
		Accounts: &UpdateMetadataAccountAccounts{
			Metadata:           &solana.AccountMeta{PublicKey: metadata, IsWritable: true},
			UpdateAuthorityKey: &solana.AccountMeta{PublicKey: updateAuthorityKey, IsSigner: true},
		},
	}
	return NewInstruction(programID, inst)
}
