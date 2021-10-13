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

type Metadata struct {
	Key                 Key
	UpdateAuthority     solana.PublicKey
	Mint                solana.PublicKey
	Data                Data
	PrimarySaleHappened bool
	IsMutable           bool
}

type Data struct {
	Name                 string
	Symbol               string
	URI                  string
	SellerFeeBasisPoints uint16
	Creators             *[]Creator `bin:"optional"`
}
type Creator struct {
	Address  solana.PublicKey
	Verified bool
	// In percentages, NOT basis points ;) Watch out!
	Share int8
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
