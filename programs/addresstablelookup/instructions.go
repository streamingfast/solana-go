package addresstablelookup

import "github.com/mr-tron/base58"

func ParseNewAddresses(addresses []byte) []string {
	var newAddresses []string
	numberOfAddresses := len(addresses) / 32

	for i := 0; i < numberOfAddresses; i++ {
		if i == numberOfAddresses {
			break
		}

		addr := addresses[(i * 32) : (i+1)*32]
		newAddresses = append(newAddresses, base58.Encode(addr))
	}

	return newAddresses
}
