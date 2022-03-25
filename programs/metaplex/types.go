package metaplex

import (
	"fmt"
	"regexp"

	"github.com/near/borsh-go"
	"github.com/streamingfast/solana-go"
)

type Key borsh.Enum

const (
	Uninitialized = iota
	EditionV1
	MasterEditionV1
	ReservationListV1
	MetadataV1
	ReservationListV2
	MasterEditionV2
	EditionMarker
)

type UseMethod borsh.Enum

const (
	Burn = iota
	Multiple
	Single
)

type TokenStandard borsh.Enum

const (
	NonFungible = iota
	FungibleAsset
	Fungible
	NonFungibleEdition
)

type Metadata struct {
	Key                 Key
	UpdateAuthority     solana.PublicKey
	Mint                solana.PublicKey
	Data                Data
	PrimarySaleHappened bool
	IsMutable           bool
	EditionNonce        *uint8         `bin:"optional"`
	TokenStandard       *TokenStandard `bin:"optional"`
	Collection          *Collection    `bin:"optional"`
	Uses                *Uses          `bin:"optional"`
}

type Data struct {
	Name                 string
	Symbol               string
	URI                  string
	SellerFeeBasisPoints uint16
	Creators             *[]Creator `bin:"optional"`
}

type DataV2 struct {
	Name                 string      `json:"name"`
	Symbol               string      `json:"symbol"`
	URI                  string      `json:"uri"`
	SellerFeeBasisPoints uint16      `json:"seller_fee_basis_points"`
	Creators             *[]Creator  `bin:"optional" json:"creators"`
	Collection           *Collection `bin:"optional"  json:"collection"`
	Uses                 *Uses       `bin:"optional"  json:"uses"`
}

type Collection struct {
	Verified bool             `json:"verified"`
	Key      solana.PublicKey `json:"key"`
}

type Uses struct {
	UseMethod UseMethod `json:"use_method"`
	Remaining uint64    `json:"remaining"`
	Total     uint64    `json:"total"`
}

type Creator struct {
	Address  solana.PublicKey `json:"address"`
	Verified bool             `json:"verified"`
	// In percentages, NOT basis points ;) Watch out!
	Share int8 `json:"share"`
}

func (m *Metadata) Decode(in []byte) error {
	err := borsh.Deserialize(m, in)
	if err != nil {
		return fmt.Errorf("unpack: %w", err)
	}
	m.Data.Name = Clean(m.Data.Name)
	m.Data.Symbol = Clean(m.Data.Symbol)
	m.Data.URI = Clean(m.Data.URI)
	return nil
}

var METADATA_REPLACE = regexp.MustCompile("\u0000")

func Clean(str string) string {
	return METADATA_REPLACE.ReplaceAllString(str, "")
}
