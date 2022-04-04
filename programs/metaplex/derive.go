package metaplex

import (
	"fmt"
	"github.com/streamingfast/solana-go"
)

func DeriveMetadataPublicKey(programID, mint solana.PublicKey) (solana.PublicKey, error) {
	path := [][]byte{
		[]byte("metadata"),
		programID[:],
		mint[:],
	}

	key, _, err := solana.PublicKeyFindProgramAddress(path, programID)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("unable to derive metaplex metadata address: %w", err)
	}
	return key, nil

}
func DeriveMetadataEditionPublicKey(programID, mint solana.PublicKey) (solana.PublicKey, error) {
	path := [][]byte{
		[]byte("metadata"),
		programID[:],
		mint[:],
		[]byte("edition"),
	}

	key, _, err := solana.PublicKeyFindProgramAddress(path, programID)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("unable to derive metaplex metadata edition address: %w", err)
	}
	return key, nil

}

func DeriveMetadataEditionCreationMarkPublicKey(programID, mint solana.PublicKey, editionNumber string) (solana.PublicKey, error) {
	path := [][]byte{
		[]byte("metadata"),
		programID[:],
		mint[:],
		[]byte("edition"),
		[]byte(editionNumber),
	}

	key, _, err := solana.PublicKeyFindProgramAddress(path, programID)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("unable to derive edition pda to mark creation: %w", err)
	}
	return key, nil

}
