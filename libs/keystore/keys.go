package keystore

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/qredo/assets/libs/crypto"
	"github.com/qredo/assets/libs/cryptowallet"
)

// Secrets - keys required for decryption and signing
type Secrets struct {
	Seed          []byte
	SikeSecretKey []byte
	BLSSecretKey  []byte
}

// GenerateBLSKeys generate BLS keys from seed
func GenerateBLSKeys(seed []byte) (blsPublic, blsSecret []byte, err error) {
	rc1, blsPublic, blsSecret := crypto.BLSKeys(seed, nil)
	if rc1 != 0 {
		err = fmt.Errorf("Failed to generate BLS keys: %v", rc1)
	}
	return
}

// GenerateSIKEKeys generate SIKE keys from seed
func GenerateSIKEKeys(seed []byte) (sikePublic, sikeSecret []byte, err error) {
	rc1, sikePublic, sikeSecret := crypto.SIKEKeys(seed)
	if rc1 != 0 {
		err = fmt.Errorf("Failed to generate SIKE keys: %v", rc1)
	}
	return
}

// GenerateECPublicKey - generate EC keys using BIP44 HD Wallets (as bitcoin) from seed
func GenerateECPublicKey(seed []byte) (ecPublic []byte, err error) {
	//EC ADD Keypair Protocol
	_, pubKeyECADD, _, err := cryptowallet.Bip44Address(seed, cryptowallet.CoinTypeBitcoinMain, 0, 0, 0)
	if err != nil {
		err = errors.Wrap(err, "Failed to derive EC HD Wallet Key")
		return
	}

	return pubKeyECADD.SerializeCompressed(), nil
}
