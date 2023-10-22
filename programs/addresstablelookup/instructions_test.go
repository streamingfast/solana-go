package addresstablelookup

import (
	"encoding/hex"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExtendLookupTableDecode(t *testing.T) {
	dataBytes, err := hex.DecodeString("02000000010000000000000000000000020000000634ff9fb638edb36d9e9e1730ef57bc9286bb905ca2a33e4eca42136e575f810277a6af97339b7ac88d1892c90446f50002309266f62e53c118244982000000")
	require.NoError(t, err)

	decodedInstruction, err := DecodeInstruction(dataBytes)
	array, err := base58.Decode("11112CCZE56WMyrQx6ToZt1KJPqyn9b3SyBFtW88us")
	var array32 [32]byte
	var expectedAddresses [][32]byte
	copy(array32[:], array)
	expectedAddresses = append(expectedAddresses, array32)

	expectedExtendLookupTableInstruction := &ExtendLookupTable{
		AddressCount: 1,
		Addresses:    expectedAddresses,
	}

	switch val := decodedInstruction.Impl.(type) {
	case *ExtendLookupTable:
		require.Equal(t, expectedExtendLookupTableInstruction, val)
	default:
		require.True(t, false)
	}
}

func TestCreateLookupTableDecode(t *testing.T) {
	dataBytes, err := hex.DecodeString("000000006697bc0900000000fd")
	require.NoError(t, err)
	decodedInstruction, err := DecodeInstruction(dataBytes)

	expectedCreateLookuptable := &CreateLookupTable{
		RecentSlot: 163354470,
		BumpSeed:   253,
	}

	switch val := decodedInstruction.Impl.(type) {
	case *CreateLookupTable:
		require.Equal(t, expectedCreateLookuptable, val)
	default:
		require.True(t, false)
	}
}
