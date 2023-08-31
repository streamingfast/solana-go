package addresstablelookup

import (
	"bytes"
)

var AddressLookupTableExtendTableInstruction = []byte{02, 00, 00, 00}

func ExtendAddressTableLookupInstruction(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	return bytes.Equal(data[:4], AddressLookupTableExtendTableInstruction)
}

func ParseNewAccounts(accounts []byte) [][]byte {
	var newAccounts [][]byte
	numberOfAccounts := len(accounts) / 32

	for i := 0; i < numberOfAccounts; i++ {
		if i == numberOfAccounts {
			break
		}

		addr := accounts[(i * 32) : (i+1)*32]
		newAccounts = append(newAccounts, addr)
	}

	return newAccounts
}
