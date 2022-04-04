package metaplex

import (
	"encoding/hex"
	"fmt"

	"go.uber.org/zap"

	"github.com/near/borsh-go"
	"github.com/streamingfast/solana-go"
)

var PROGRAM_ID = solana.MustPublicKeyFromBase58("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")

type InstType uint8

const (
	CreateMetadataAccountV1Inst InstType = iota
	UpdateMetadataAccountV1Inst
	DeprecatedCreateMasterEditionInst
	DeprecatedMintNewEditionFromMasterEditionViaPrintingTokenInst
	UpdatePrimarySaleHappenedViaTokenInst
	DeprecatedSetReservationListInst
	DeprecatedCreateReservationListInst
	SignMetadataInst
	DeprecatedMintPrintingTokensViaTokenInst
	DeprecatedMintPrintingTokensInst
	CreateMasterEditionInst
	MintNewEditionFromMasterEditionViaTokenInst
	ConvertMasterEditionV1ToV2Inst
	MintNewEditionFromMasterEditionViaVaultProxyInst
	PuffMetadataInst
	UpdateMetadataAccountV2Inst
	CreateMetadataAccountV2Inst
	CreateMasterEditionV3Inst
	VerifyCollectionInst
	UtilizeInst
	ApproveUseAuthorityInst
	RevokeUseAuthorityInst
	UnverifyCollectionInst
	ApproveCollectionAuthorityInst
	RevokeCollectionAuthorityInst
)

type AccountListable interface {
	ListAccounts() []*solana.AccountMeta
}
type Instruction struct {
	Impl interface{}

	programId solana.PublicKey
}

func NewInstruction(programId solana.PublicKey, impl interface{}) *Instruction {
	return &Instruction{
		Impl:      impl,
		programId: programId,
	}
}

func (i *Instruction) Data() ([]byte, error) {
	data, err := borsh.Serialize(i.Impl)
	if err != nil {
		return nil, fmt.Errorf("borsh serailize: %w", err)
	}
	zlog.Debug("encodded metaplex instruction", zap.String("data", hex.EncodeToString(data)))
	return data, nil
}

func (i *Instruction) Accounts() (out []*solana.AccountMeta) {
	if listeable, ok := i.Impl.(AccountListable); ok {
		return listeable.ListAccounts()
	}
	panic("an instruction needs to implement the Accounts()")
}

func (i *Instruction) ProgramID() solana.PublicKey {
	return i.programId
}
