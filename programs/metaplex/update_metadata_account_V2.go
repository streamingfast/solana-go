package metaplex

import (
	"github.com/streamingfast/solana-go"
)

type UpdateMetadataV2Account struct {
	Instruction         InstType
	Data                *DataV2
	UpdateAuthority     *solana.PublicKey
	PrimarySaleHappened *bool
	IsMutable           *bool

	Accounts *UpdateMetadataAccountV2Accounts `borsh_skip:"true"`
}

func (i UpdateMetadataV2Account) ListAccounts() []*solana.AccountMeta {
	return []*solana.AccountMeta{
		i.Accounts.Metadata,
		i.Accounts.UpdateAuthorityKey,
	}
}

type UpdateMetadataAccountV2Accounts struct {
	Metadata           *solana.AccountMeta
	UpdateAuthorityKey *solana.AccountMeta
}

func NewUpdateMetadataAccountV2Instruction(
	programID solana.PublicKey,
	data *DataV2,
	updateAuthority *solana.PublicKey,
	primarySaleHappened *bool,
	isMutable *bool,
	metadata solana.PublicKey,
	updateAuthorityKey solana.PublicKey,
) *Instruction {
	var inst = UpdateMetadataV2Account{
		Instruction:         UpdateMetadataAccountV2Inst,
		Data:                data,
		UpdateAuthority:     updateAuthority,
		PrimarySaleHappened: primarySaleHappened,
		IsMutable:           isMutable,
		Accounts: &UpdateMetadataAccountV2Accounts{
			Metadata:           &solana.AccountMeta{PublicKey: metadata, IsWritable: true},
			UpdateAuthorityKey: &solana.AccountMeta{PublicKey: updateAuthorityKey, IsSigner: true},
		},
	}
	return NewInstruction(programID, inst)
}
