package metaplex

import (
	"github.com/streamingfast/solana-go"
)

type UpdateMetadataV1Account struct {
	Instruction         InstType
	Data                *Data
	UpdateAuthority     *solana.PublicKey
	PrimarySaleHappened *bool

	Accounts *UpdateMetadataAccountV1Accounts `borsh_skip:"true"`
}

func (i UpdateMetadataV1Account) ListAccounts() []*solana.AccountMeta {
	return []*solana.AccountMeta{
		i.Accounts.Metadata,
		i.Accounts.UpdateAuthorityKey,
	}
}

/// Update a Metadata
///   0. `[writable]` Metadata account
///   1. `[signer]` Update authority key
type UpdateMetadataAccountV1Accounts struct {
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

func NewUpdateMetadataAccountV1Instruction(
	programID solana.PublicKey,
	data *Data,
	updateAuthority *solana.PublicKey,
	primarySaleHappened *bool,
	metadata solana.PublicKey,
	updateAuthorityKey solana.PublicKey,
) *Instruction {
	var inst = UpdateMetadataV1Account{
		Instruction:         UpdateMetadataAccountV1Inst,
		Data:                data,
		UpdateAuthority:     updateAuthority,
		PrimarySaleHappened: primarySaleHappened,
		Accounts: &UpdateMetadataAccountV1Accounts{
			Metadata:           &solana.AccountMeta{PublicKey: metadata, IsWritable: true},
			UpdateAuthorityKey: &solana.AccountMeta{PublicKey: updateAuthorityKey, IsSigner: true},
		},
	}
	return NewInstruction(programID, inst)
}
