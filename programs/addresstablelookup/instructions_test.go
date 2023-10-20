package addresstablelookup

import (
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ParseNewAddresses(t *testing.T) {
	tests := []struct {
		name             string
		input            []byte
		expectedAccounts []string
	}{
		{
			name: "multiple accounts",
			input: []byte{
				237, 231, 230, 250, 137, 55, 19, 37, 120, 93, 19, 216, 54, 107, 151, 114, 155, 24, 141, 211, 190, 191, 218, 243, 87, 87, 54, 115, 166, 12, 188, 22, 162, 139, 100, 91, 179, 213, 6, 44, 73, 187, 105, 43, 1, 170, 52, 130, 216, 82, 234, 12, 95, 37, 109, 124, 179, 40, 191, 196, 220, 149, 5, 70, 225, 202, 101, 93, 88, 45, 114, 49, 60, 190, 66, 114, 194, 187, 188, 16, 27, 249, 46, 41, 245, 210, 142, 237, 201, 247, 189, 71, 162, 120, 59, 231, 126, 84, 119, 26, 87, 166, 241, 76, 169, 228, 2, 213, 74, 238, 69, 247, 55, 138, 202, 54, 92, 123, 22, 154, 126, 200, 63, 81, 130, 178, 152, 240, 138, 239, 125, 247, 227, 166, 116, 164, 62, 209, 129, 86, 34, 14, 3, 6, 250, 71, 135, 139, 147, 32, 126, 36, 149, 229, 158, 237, 187, 218, 76, 161,
			},
			expectedAccounts: []string{
				"H1gikkvnijbQeGLi3tk7RgY4nhiTkqQAmyHX6fBduwbj",
				"BwWKuqJyKVTYyaQJzSouQQ9pbuocNasYH3dK9RDRQEbf",
				"GCPjTU3KDUyHUtHh2Z6d1fvHvQztGn9T1yif7ENEAYgW",
				"9W959DqEETiGZocYWCQPaJ6sBmUzgfxXfqGeTEdp3aQP",
				"AMM55ShdkoGRB5jVYPjWziwk8m5MpwyDgsMWHaMSQWH6",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			newAccounts := ParseNewAccounts(test.input)
			for i, account := range newAccounts {
				require.Equal(t, test.expectedAccounts[i], base58.Encode(account))
			}
		})
	}
}

func Test_ExtendAddressTableLookupInstruction(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected bool
	}{
		{
			name:     "valid instruction",
			input:    []byte{02, 00, 00, 00},
			expected: true,
		},
		{
			name:     "invalid instruction",
			input:    []byte{02, 00, 00, 01},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expected, ExtendAddressTableLookupInstruction(test.input))
		})
	}
}

func TestExtendLookupTable(t *testing.T) {

}
