package metaplex

import (
	"github.com/streamingfast/solana-go"
	"github.com/streamingfast/solana-go/programs/system"
	"github.com/streamingfast/solana-go/programs/token"
)

type MintNewEditionFromMasterEditionViaToken struct {
	Instruction InstType
	Edition     uint64

	Accounts *MintNewEditionFromMasterEditionViaTokenAccounts `borsh_skip:"true"`
}

//
func (i MintNewEditionFromMasterEditionViaToken) ListAccounts() []*solana.AccountMeta {
	acounts := []*solana.AccountMeta{
		i.Accounts.NewMetadata,
		i.Accounts.NewEdition,
		i.Accounts.MasterRecordEdition,
		i.Accounts.Mint,
		i.Accounts.Edition,
		i.Accounts.MintAuthority,
		i.Accounts.Payer,
		i.Accounts.Owner,
		i.Accounts.TokenAccount,
		i.Accounts.UpdateAuthority,
		i.Accounts.MasterRecordMetadata,
		i.Accounts.TokenProgram,
		i.Accounts.SystemProgram,
		i.Accounts.RentProgram,
	}
	return acounts

}

/// Given a token account containing the master edition token to prove authority, and a brand new non-metadata-ed mint with one token
/// make a new Metadata + Edition that is a child of the master edition denoted by this authority token.
///   0. `[writable]` New Metadata key (pda of ['metadata', program id, mint id])
///   1. `[writable]` New Edition (pda of ['metadata', program id, mint id, 'edition'])
///   2. `[writable]` Master Record Edition V2 (pda of ['metadata', program id, master metadata mint id, 'edition'])
///   3. `[writable]` Mint of new token - THIS WILL TRANSFER AUTHORITY AWAY FROM THIS KEY
///   4. `[writable]` Edition pda to mark creation - will be checked for pre-existence. (pda of ['metadata', program id, master metadata mint id, 'edition', edition_number])
///   where edition_number is NOT the edition number you pass in args but actually edition_number = floor(edition/EDITION_MARKER_BIT_SIZE).
///   5. `[signer]` Mint authority of new mint
///   6. `[signer]` payer
///   7. `[signer]` owner of token account containing master token (#8)
///   8. `[]` token account containing token from master metadata mint
///   9. `[]` Update authority info for new metadata
///   10. `[]` Master record metadata account
///   11. `[]` Token program
///   12. `[]` System program
///   13. `[]` Rent info
type MintNewEditionFromMasterEditionViaTokenAccounts struct {
	NewMetadata          *solana.AccountMeta
	NewEdition           *solana.AccountMeta
	MasterRecordEdition  *solana.AccountMeta
	Mint                 *solana.AccountMeta
	Edition              *solana.AccountMeta
	MintAuthority        *solana.AccountMeta
	Payer                *solana.AccountMeta
	Owner                *solana.AccountMeta
	TokenAccount         *solana.AccountMeta
	UpdateAuthority      *solana.AccountMeta
	MasterRecordMetadata *solana.AccountMeta
	TokenProgram         *solana.AccountMeta
	SystemProgram        *solana.AccountMeta
	RentProgram          *solana.AccountMeta
}

func NewMintNewEditionFromMasterEditionViaToken(
	programID solana.PublicKey,
	editionNum uint64,
	newMetadata,
	newEdition,
	masterRecordEdition,
	mint,
	edition,
	mintAuthority,
	payer,
	owner,
	tokenAccount,
	updateAuthority,
	masterRecordMetadata solana.PublicKey,
) *Instruction {
	var inst = MintNewEditionFromMasterEditionViaToken{
		Instruction: MintNewEditionFromMasterEditionViaTokenInst,
		Edition:     editionNum,
		Accounts: &MintNewEditionFromMasterEditionViaTokenAccounts{
			NewMetadata:          &solana.AccountMeta{PublicKey: newMetadata, IsWritable: true},
			NewEdition:           &solana.AccountMeta{PublicKey: newEdition, IsWritable: true},
			MasterRecordEdition:  &solana.AccountMeta{PublicKey: masterRecordEdition, IsWritable: true},
			Mint:                 &solana.AccountMeta{PublicKey: mint, IsWritable: true},
			Edition:              &solana.AccountMeta{PublicKey: edition, IsWritable: true},
			MintAuthority:        &solana.AccountMeta{PublicKey: mintAuthority, IsSigner: true},
			Payer:                &solana.AccountMeta{PublicKey: payer, IsSigner: true},
			Owner:                &solana.AccountMeta{PublicKey: owner, IsSigner: true},
			TokenAccount:         &solana.AccountMeta{PublicKey: tokenAccount},
			UpdateAuthority:      &solana.AccountMeta{PublicKey: updateAuthority},
			MasterRecordMetadata: &solana.AccountMeta{PublicKey: masterRecordMetadata},
			TokenProgram:         &solana.AccountMeta{PublicKey: token.PROGRAM_ID},
			SystemProgram:        &solana.AccountMeta{PublicKey: system.PROGRAM_ID},
			RentProgram:          &solana.AccountMeta{PublicKey: system.SYSVAR_RENT},
		},
	}
	return NewInstruction(programID, inst)
}
