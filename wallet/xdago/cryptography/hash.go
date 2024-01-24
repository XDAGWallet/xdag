package cryptography

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"goXdagWallet/xdago/common"
	"goXdagWallet/xdago/secp256k1"

	"golang.org/x/crypto/ripemd160"
)

func HashTwice(input []byte) [32]byte {
	h := sha256.Sum256(input)
	return sha256.Sum256(h[:])
}

func Sha256Hash160(input []byte) (out common.Hash160) {
	h := sha256.Sum256(input)
	rim := ripemd160.New()
	rim.Write(h[:])
	copy(out[:], rim.Sum(nil))
	return
}

func HmacSha512(key, input []byte) []byte {
	mac := hmac.New(sha512.New, key)
	mac.Write(input)
	return mac.Sum(nil)
}

func ToBytesAddress(key *secp256k1.PrivateKey) common.Hash160 {
	return Sha256Hash160(key.PubKey().SerializeCompressed())
}
