package utils

import (
	"errors"
	"goXdagWallet/xdago/common"
	"math"
	"unsafe"
)

const (
	bits2mime = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

var (
	mime2bits = make([]byte, 256)

	// 8 bytes with zero data
	Zero8bytes = make([]byte, 8)
)

func init() {
	for i := range mime2bits {
		mime2bits[i] = 0xFF
	}
	var i uint8
	for i = 0; i < 64; i++ {
		mime2bits[bits2mime[i]] = i
	}
}

// Hash2Address converts hash to address
func Hash2Address(h common.Hash) string {
	address := make([]byte, common.XDAG_ADDRESS_SIZE)
	var c, d, j uint
	// every 3 bytes(24 bits) hashs convert to 4 chars(6 bit each)
	// first 24 bytes hash to 32 byte address, ignore last 8 bytes of hash
	for i := 0; i < common.XDAG_ADDRESS_SIZE; i++ {
		if d < 6 {
			d += 8
			c <<= 8
			c |= uint(h[j])
			j++
		}
		d -= 6
		address[i] = bits2mime[c>>d&0x3F]
	}
	return bytes2str(address)
}

// Address2Hash converts address to hash
func Address2Hash(addr string) (common.Hash, error) {
	var hash common.Hash
	var i, e, n, j uint
	var c, d uint8
	if len(addr) != 32 {
		return hash, errors.New("address length error")
	}
	// convert 32 byte address to 24 bytes hash
	// each byte (8 bits) address to 6 bits hash
	for i = 0; i < common.XDAG_ADDRESS_SIZE; i++ {
		for {
			c = addr[i]
			if c == 0 {
				return hash, errors.New("address string error")
			}
			d = mime2bits[c]
			if d&0xC0 == 0 {
				break
			}
		}
		e <<= 6
		e |= uint(d)
		n += 6
		if n >= 8 {
			n -= 8
			hash[j] = uint8(e >> n)
			j++
		}
	}
	copy(hash[24:], Zero8bytes) // set last 8 bytes of hash to 0
	return hash, nil
}

func Xdag2Amount(value float64) uint64 {
	amount := uint64(math.Floor(value))
	res := amount << 32
	value = value - float64(amount)
	value = value * math.Pow(2, 32)
	amount = uint64(math.Ceil(value))
	return res + amount
}

//// RawBlock contains raw XDAG block bytes
//type RawBlock struct {
//	Hash      [32]byte
//	Address   string
//	Timestamp uint64
//	RawBytes  []byte
//}
//
//// NewRawBlock builds new raw block from bytes
//func NewRawBlock(b []byte) RawBlock {
//
//	header := make([]byte, 8)
//	copy(header, b[:8])     // backup block transport header
//	copy(b[:8], Zero8bytes) // clear block transport header
//
//	hash := sha256.Sum256(b)
//	copy(b[:8], header) // restore block transport header
//	r := RawBlock{
//		Hash:     sha256.Sum256(hash[:]),
//		RawBytes: b,
//	}
//	// get time from block header
//	r.Timestamp = binary.LittleEndian.Uint64(b[16:24])
//
//	r.Address = Hash2Address(r.Hash)
//	return r
//}

// unsafe and fast convert string to bytes slice
func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// unsafe and fast convert bytes slice to string
func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
