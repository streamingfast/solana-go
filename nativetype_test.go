package solana

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMustPublicKeyFromBase58(t *testing.T) {
	require.Panics(t, func() {
		MustPublicKeyFromBase58("toto")
	})
}

func TestSignature_Verify(t *testing.T) {
	type input struct {
		publickKey string
		message    string
		signature  string
	}

	tests := []struct {
		name     string
		in       input
		expected bool
	}{
		{
			"pass",
			input{
				publickKey: "Gkodg1n3z56G8XWohmNSwPNYYtsfrc1AFCZndT7XjjZ2",
				message:    "Signing a message that will prove to the server you own public key associated to your wallet",
				signature:  "7003d58a9db0ad2cb288014bf41ad60f4ad8d207b232b4adbb37933850e7c96693547e63bf977ed98f67e7b6911c376236ac14095c57b984afc14dc2c1bec308",
			},
			true,
		},

		{
			"fails",
			input{
				publickKey: "Gkodg1n3z56G8XWohmNSwPNYYtsfrc1AFCZndT7XjjZ2",
				message:    "Signing a message that will prove to the server you own public key associated to your wallet",
				signature:  "8003d58a9db0ad2cb288014bf41ad60f4ad8d207b232b4adbb37933850e7c96693547e63bf977ed98f67e7b6911c376236ac14095c57b984afc14dc2c1bec308",
			},
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			publicKey := MustPublicKeyFromBase58(test.in.publickKey)
			message := []byte(test.in.message)
			signature, err := NewSignatureFromString(test.in.signature)
			require.NoError(t, err)

			if test.expected {
				require.True(t, signature.Verify(publicKey, message), "Signature %s is invalid for public key %s (message %q)", test.in.signature, publicKey, test.in.message)
			} else {
				require.False(t, signature.Verify(publicKey, message), "Signature %s is valid but should not have been for public key %s (message %s)", test.in.signature, publicKey, test.in.message)
			}
		})
	}
}
