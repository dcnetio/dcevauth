package util

import (
	"crypto/ed25519"

	"github.com/ChainSafe/go-schnorrkel"
	"github.com/libp2p/go-libp2p/core/crypto"
)

func ImportTestMnemonic(mnemonic string) (prikey crypto.PrivKey, err error) {
	seed, err := schnorrkel.SeedFromMnemonic(mnemonic, "")
	if err != nil {
		return
	}

	secret := ed25519.NewKeyFromSeed(seed[:32])
	prikey, err = crypto.UnmarshalEd25519PrivateKey([]byte(secret))
	if err != nil {
		return
	}
	return
}
