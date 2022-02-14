package token

import (
	"encoding/hex"
	"testing"

	"github.com/streamingfast/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSetAuthorityInstruction(t *testing.T) {
	account := solana.MustPublicKeyFromBase58("2kPGkTUzGZDTwxamiC99vggZZ52Dj6TKLNTErXmbNVwt")
	newAuthority := solana.MustPublicKeyFromBase58("Gr5UanqwiKA54GGnw4b1bB5M8eatQzj6s6FQ9FeTze5C")
	authorityType := AccountOwnerAuthorityType
	currentAuthority := solana.MustPublicKeyFromBase58("Gg1CWowuNc9ytKuX1p7hZ2mgBmWrjT9eNeXf8gvdq1Kd")

	inst := NewSetAuthorityInstruction(
		account,
		newAuthority,
		authorityType,
		currentAuthority,
	)

	cnt, err := inst.Data()
	require.NoError(t, err)
	assert.Equal(t, "060201eb71d2f68370da513ad471ed6bb975a331931cb4e221aead8247ab8e4618d04f", hex.EncodeToString(cnt))
}
