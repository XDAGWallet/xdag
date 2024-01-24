package cryptography

import (
	"crypto/sha256"
	"encoding/hex"
	"goXdagWallet/xdago/secp256k1"
	"testing"
)

func TestEcdsaSign(t *testing.T) {
	pkBytes, err := hex.DecodeString("22a47fa09a223f2aa079edf85a7c2d4f87" +
		"20ee63e502ee2869afab7de234b80c")
	if err != nil {
		panic(err)
	}
	privKey := secp256k1.PrivKeyFromBytes(pkBytes)
	message := "test message"
	hash := sha256.Sum256([]byte(message))
	r, s := EcdsaSign(privKey, hash[:])

	pubKey := privKey.PubKey()

	ok := EcdsaVerify(pubKey, hash[:], r[:], s[:])
	if !ok {
		panic("verify failed")
	}

}
