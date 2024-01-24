// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2015-2021 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package ecdsa_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"goXdagWallet/xdago/secp256k1"
	"goXdagWallet/xdago/secp256k1/ecdsa"
)

// This example demonstrates signing a message with a secp256k1 private key that
// is first parsed from raw bytes and serializing the generated signature.
func ExampleSign() {
	// Decode a hex-encoded private key.
	pkBytes, err := hex.DecodeString("22a47fa09a223f2aa079edf85a7c2d4f87" +
		"20ee63e502ee2869afab7de234b80c")
	if err != nil {
		fmt.Println(err)
		return
	}
	privKey := secp256k1.PrivKeyFromBytes(pkBytes)

	// Sign a message using the private key.
	message := "test message"
	messageHash := sha256.Sum256([]byte(message))
	signature := ecdsa.Sign(privKey, messageHash[:])

	// Serialize and display the signature.
	fmt.Printf("Serialized Signature: %x\n", signature.Serialize())

	// Verify the signature for the message using the public key.
	pubKey := privKey.PubKey()
	verified := signature.Verify(messageHash[:], pubKey)
	fmt.Printf("Signature Verified? %v\n", verified)

	// Output:
	// Serialized Signature: 3045022100fa7cdbd9243b99889b033e88ae2ddf55cc189efd5ae64dfa77655f01fc48e8000220045ec2f0dfebc7891d31b40d1ed686ca0e33c7c1b1b693e0fb305e6fc4d84a6a
	// Signature Verified? true
}

// This example demonstrates verifying a secp256k1 signature against a public
// key that is first parsed from raw bytes.  The signature is also parsed from
// raw bytes.
func ExampleSignature_Verify() {
	// Decode hex-encoded serialized public key.
	pubKeyBytes, err := hex.DecodeString("02a673638cb9587cb68ea08dbef685c" +
		"6f2d2a751a8b3c6f2a7e9a4999e6e4bfaf5")
	if err != nil {
		fmt.Println(err)
		return
	}
	pubKey, err := secp256k1.ParsePubKey(pubKeyBytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Decode hex-encoded serialized signature.
	sigBytes, err := hex.DecodeString("3045022100fa7cdbd9243b99889b033e88ae" +
		"2ddf55cc189efd5ae64dfa77655f01fc48e8000220045ec2f0dfebc7891d31b40d1" +
		"ed686ca0e33c7c1b1b693e0fb305e6fc4d84a6a")
	if err != nil {
		fmt.Println(err)
		return
	}
	signature, err := ecdsa.ParseDERSignature(sigBytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Verify the signature for the message using the public key.
	message := "test message"
	messageHash := sha256.Sum256([]byte(message))
	verified := signature.Verify(messageHash[:], pubKey)
	fmt.Println("Signature Verified?", verified)

	// Output:
	// Signature Verified? true
}
