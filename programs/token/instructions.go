package token

import (
	"bytes"
	"fmt"

	bin "github.com/streamingfast/binary"
	"github.com/streamingfast/solana-go"
	"github.com/streamingfast/solana-go/text"
)

var PROGRAM_ID = solana.MustPublicKeyFromBase58("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")

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
	{"initialize_mint", (*InitializeMint)(nil)},
	{"initialize_account", (*InitializeAccount)(nil)},
	{"InitializeMultisig", (*InitializeMultisig)(nil)},
	{"Transfer", (*Transfer)(nil)},
	{"Approve", (*Approve)(nil)},
	{"Revoke", (*Revoke)(nil)},
	{"SetAuthority", (*SetAuthority)(nil)},
	{"MintTo", (*MintTo)(nil)},
	{"Burn", (*Burn)(nil)},
	{"CloseAccount", (*CloseAccount)(nil)},
	{"FreezeAccount", (*FreezeAccount)(nil)},
	{"ThawAccount", (*ThawAccount)(nil)},
	{"TransferChecked", (*TransferChecked)(nil)},
	{"ApproveChecked", (*ApproveChecked)(nil)},
	{"MintToChecked", (*MintToChecked)(nil)},
	{"BurnChecked", (*BurnChecked)(nil)},
})

type Instruction struct {
	bin.BaseVariant
}

func (i *Instruction) Accounts() (out []*solana.AccountMeta) {
	switch i.TypeID {
	case 0:
		accounts := i.Impl.(*InitializeMint).Accounts
		out = []*solana.AccountMeta{accounts.Mint, accounts.RentProgram}
	case 3:
		accounts := i.Impl.(*Transfer).Accounts
		out = []*solana.AccountMeta{
			accounts.Source,
			accounts.Destination,
			accounts.From,
		}
	case 6:
		accounts := i.Impl.(*SetAuthority).Accounts
		out = []*solana.AccountMeta{
			accounts.Account,
			accounts.CurrentAuthority,
		}
	case 7:
		accounts := i.Impl.(*MintTo).Accounts
		out = []*solana.AccountMeta{
			accounts.Mint,
			accounts.Account,
			accounts.MintAuthority,
		}
	case 9:
		accounts := i.Impl.(*CloseAccount).Accounts
		out = []*solana.AccountMeta{
			accounts.Account,
			accounts.Destination,
			accounts.Owner,
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

type InitializeMultisigAccounts struct {
}
type InitializeMultisig struct {
	Accounts *InitializeMultisigAccounts
}

type InitializeMintAccounts struct {
	Mint        *solana.AccountMeta
	RentProgram *solana.AccountMeta
	///   0. `[writable]` The mint to initialize.
	///   1. `[]` Rent sysvar
}
type InitializeMint struct {
	/// Number of base 10 digits to the right of the decimal place.
	Decimals uint8
	/// The authority/multisignature to mint tokens.
	MintAuthority solana.PublicKey
	/// The freeze authority/multisignature of the mint.
	FreezeAuthority *solana.PublicKey       `bin:"optional"`
	Accounts        *InitializeMintAccounts `bin:"-"`
}

func NewInitializeMintInstruction(
	decimals uint8,
	mint solana.PublicKey,
	mintAuthority solana.PublicKey,
	freezeAuthority *solana.PublicKey,
	rentProgram solana.PublicKey,
) *Instruction {
	return &Instruction{
		BaseVariant: bin.BaseVariant{
			TypeID: 0,
			Impl: &InitializeMint{
				Decimals:        decimals,
				MintAuthority:   mintAuthority,
				FreezeAuthority: freezeAuthority,
				Accounts: &InitializeMintAccounts{
					Mint:        &solana.AccountMeta{PublicKey: mint, IsWritable: true},
					RentProgram: &solana.AccountMeta{PublicKey: rentProgram},
				},
			},
		},
	}
}

func (i *InitializeMint) SetAccounts(accounts []*solana.AccountMeta) error {
	i.Accounts = &InitializeMintAccounts{
		Mint:        accounts[0],
		RentProgram: accounts[1],
	}
	return nil
}

type TransferAccounts struct {
	Source      *solana.AccountMeta
	Destination *solana.AccountMeta
	From        *solana.AccountMeta
}

type Transfer struct {
	Amount   uint64
	Accounts *TransferAccounts
}

func NewTransferInstruction(
	amount uint64,
	source solana.PublicKey,
	destination solana.PublicKey,
	from solana.PublicKey,
) *Instruction {
	return &Instruction{
		BaseVariant: bin.BaseVariant{
			TypeID: 3,
			Impl: &Transfer{
				Amount: amount,
				Accounts: &TransferAccounts{
					Source:      &solana.AccountMeta{PublicKey: source, IsWritable: true},
					Destination: &solana.AccountMeta{PublicKey: destination, IsWritable: true},
					From:        &solana.AccountMeta{PublicKey: from, IsSigner: true},
				},
			},
		},
	}
}

type ApproveAccounts struct {
}
type Approve struct {
	Accounts *ApproveAccounts
}

type RevokeAccounts struct {
}
type Revoke struct {
	Accounts *RevokeAccounts
}

type SetAuthorityAccounts struct {
	Account          *solana.AccountMeta
	CurrentAuthority *solana.AccountMeta
}
type AuthorityType byte

const (
	MintTokensAuthorityType AuthorityType = iota
	FreezeAccountAuthorityType
	AccountOwnerAuthorityType
	CloseAccountAuthorityType
)

type SetAuthority struct {
	AuthorityType   AuthorityType
	NewAuthorityKey solana.PublicKey      `bin:"optional"`
	Accounts        *SetAuthorityAccounts `bin:"-"`
}

func NewSetAuthorityInstruction(
	account solana.PublicKey, // either a spl_token::mint or spl_token::account
	newAuthority solana.PublicKey,
	authorityType AuthorityType,
	currentAuthority solana.PublicKey,
) *Instruction {
	return &Instruction{
		BaseVariant: bin.BaseVariant{
			TypeID: 6,
			Impl: &SetAuthority{
				AuthorityType:   authorityType,
				NewAuthorityKey: newAuthority,
				Accounts: &SetAuthorityAccounts{
					Account:          &solana.AccountMeta{PublicKey: account, IsWritable: true},
					CurrentAuthority: &solana.AccountMeta{PublicKey: currentAuthority, IsSigner: true},
				},
			},
		},
	}
}

type MintToAccounts struct {
	Mint          *solana.AccountMeta
	Account       *solana.AccountMeta
	MintAuthority *solana.AccountMeta

	///   0. `[writable]` The mint.
	///   1. `[writable]` The account to mint tokens to.
	///   2. `[signer]` The mint's minting authority.

}
type MintTo struct {
	Amount   uint64
	Accounts *MintToAccounts `bin:"-"`
}

func NewMintTo(
	Amount uint64,
	Mint solana.PublicKey,
	Account solana.PublicKey,
	MintAuthority solana.PublicKey,
) *Instruction {
	return &Instruction{
		BaseVariant: bin.BaseVariant{
			TypeID: 7,
			Impl: &MintTo{
				Amount: Amount,
				Accounts: &MintToAccounts{
					Mint:          &solana.AccountMeta{PublicKey: Mint, IsWritable: true},
					Account:       &solana.AccountMeta{PublicKey: Account, IsWritable: true},
					MintAuthority: &solana.AccountMeta{PublicKey: MintAuthority, IsSigner: true},
				},
			},
		},
	}
}

type BurnAccounts struct {
}
type Burn struct {
	Accounts *BurnAccounts
}

type CloseAccountAccounts struct {
	Account     *solana.AccountMeta
	Destination *solana.AccountMeta
	Owner       *solana.AccountMeta
}

type CloseAccount struct {
	Accounts *CloseAccountAccounts
}

func NewCloseAccount(
	account solana.PublicKey,
	destination solana.PublicKey,
	owner solana.PublicKey,
) *Instruction {
	return &Instruction{
		BaseVariant: bin.BaseVariant{
			TypeID: 9,
			Impl: &CloseAccount{
				Accounts: &CloseAccountAccounts{
					Account:     &solana.AccountMeta{PublicKey: account, IsWritable: true},
					Destination: &solana.AccountMeta{PublicKey: destination, IsWritable: true},
					Owner:       &solana.AccountMeta{PublicKey: owner, IsSigner: true},
				},
			},
		},
	}
}

type FreezeAccountAccounts struct {
}
type FreezeAccount struct {
	Accounts *FreezeAccountAccounts
}

type ThawAccountAccounts struct {
}
type ThawAccount struct {
	Accounts *ThawAccountAccounts
}

type TransferCheckedAccounts struct {
}
type TransferChecked struct {
	Accounts *TransferCheckedAccounts
}

type ApproveCheckedAccounts struct {
}
type ApproveChecked struct {
	Accounts *ApproveCheckedAccounts
}

type MintToCheckedAccounts struct {
}
type MintToChecked struct {
	Accounts *MintToCheckedAccounts
}

type BurnCheckedAccounts struct {
}
type BurnChecked struct {
	Accounts *BurnCheckedAccounts
}

type InitializeAccountAccounts struct {
	Account    *solana.AccountMeta `text:"linear,notype"`
	Mint       *solana.AccountMeta `text:"linear,notype"`
	Owner      *solana.AccountMeta `text:"linear,notype"`
	RentSysvar *solana.AccountMeta `text:"linear,notype"`
}

type InitializeAccount struct {
	Accounts *InitializeAccountAccounts `bin:"-"`
}

func (i *InitializeAccount) SetAccounts(accounts []*solana.AccountMeta) error {
	i.Accounts = &InitializeAccountAccounts{
		Account:    accounts[0],
		Mint:       accounts[1],
		Owner:      accounts[2],
		RentSysvar: accounts[3],
	}
	return nil
}
