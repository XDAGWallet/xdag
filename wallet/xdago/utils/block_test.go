package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"goXdagWallet/xdago/base58"
	"goXdagWallet/xdago/common"
	"testing"
)

func TestAddress(t *testing.T) {
	bs := "cd999e172c0d3e36850b6a3d5202638e4735bfc83be0ca0b"
	b, _ := hex.DecodeString(bs)
	var h common.Hash
	copy(h[:], b[:])
	a := Hash2Address(h)
	fmt.Println(a)

	h2, err := Address2Hash(a)
	if err != nil {
		panic(err)
	}

	fmt.Println(h2)

	addrBytes, _, err := base58.ChkDec(a)
	if err != nil {
		panic(err)
	}
	if len(addrBytes) != 24 {
		panic(errors.New("transaction receive address length error"))
	}
	fmt.Println(addrBytes)

}
