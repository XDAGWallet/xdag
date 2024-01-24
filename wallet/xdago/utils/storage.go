package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"goXdagWallet/config"
	"goXdagWallet/xdago/common"
	"goXdagWallet/xdago/cryptography"
	"io"
	"os"
	"path"
)

const (
	xdagjDatFolder   = "xdagj_dat"
	xdagStoreFolder  = "storage"
	xdagStoreTestNet = "storage-testnet"
	storeFileExt     = ".dat"
)

func makeDir1(t uint64) string {
	var dir string
	if config.GetConfig().Option.IsTestNet {
		dir = xdagStoreTestNet
	} else {
		dir = xdagStoreFolder
	}
	subdir := fmt.Sprintf("%02x", uint8(t>>40))

	return path.Join(xdagjDatFolder, dir, subdir)
}

func makeDir2(t uint64) string {

	dir := makeDir1(t)
	subdir := fmt.Sprintf("%02x", uint8(t>>32))

	return path.Join(dir, subdir)
}

func makeDir3(t uint64) string {

	dir := makeDir2(t)
	subdir := fmt.Sprintf("%02x", uint8(t>>24))

	return path.Join(dir, subdir)
}

func makeFile(t uint64) string {

	dir := makeDir3(t)
	subdir := fmt.Sprintf("%02x", uint8(t>>16))

	return path.Join(dir, subdir) + storeFileExt
}

// LoadBlock loads first wallet block from XDAG storage, ignore check sum
func LoadBlock(startTime, endTime uint64) ([]string, error) {
	var mask uint64
	var addrsses []string
	var hash [32]byte
	for startTime < endTime {
		datPath := makeFile(startTime)
		//fmt.Println(datPath)
		file, err := os.OpenFile(datPath, os.O_RDONLY, os.ModePerm)
		if file != nil && err == nil {
			fileInfo, err := file.Stat()
			if err != nil {
				return nil, err
			}
			n := fileInfo.Size()                         // file size
			if n%common.XDAG_BLOCK_SIZE != 0 || n == 0 { // n should be integral multiple of block size
				return nil, errors.New("file size error")
			}
			var buffer bytes.Buffer
			_, err = io.CopyN(&buffer, file, common.XDAG_BLOCK_SIZE) // read a block
			if err != nil {
				// not EOF
				return nil, err
			}
			mask = (uint64(1) << 16) - 1
			file.Close()
			block := buffer.Bytes()
			fieldTypes := binary.LittleEndian.Uint64(block[8:16])
			// header(1/8),5(sign_r),5(sign_s)
			if fieldTypes == 0x0551 || fieldTypes == 0x0558 {
				// return block, nil
				hash = cryptography.HashTwice(block)
				addrsses = append(addrsses, Hash2Address(hash))
			}
			// else {
			// 	return nil, errors.New("block type error")
			// }
		} else if FileExists(makeDir3(startTime)) {
			mask = (uint64(1) << 16) - 1
		} else if FileExists(makeDir2(startTime)) {
			mask = (uint64(1) << 24) - 1
		} else if FileExists(makeDir1(startTime)) {
			mask = (uint64(1) << 32) - 1
		} else {
			mask = (uint64(1) << 40) - 1
		}
		startTime |= mask
		startTime++
	}
	// return nil, errors.New("load block error")
	return addrsses, nil
}

func AddressFromStorage() ([]string, error) {
	var begin uint64
	if config.GetConfig().Option.IsTestNet {
		begin = XDAG_TEST_ERA
	} else {
		begin = XDAG_MAIN_ERA
	}
	// var res []byte
	// block, err := LoadBlock(begin, GetCurrentTimestamp())
	addr, err := LoadBlock(begin, GetCurrentTimestamp())
	if err != nil {
		return nil, err
	}
	if len(addr) < 1 {
		return nil, errors.New("wallet address not found")
	}
	return addr, nil
	// for err == nil {
	// 	res = block
	// 	t := binary.LittleEndian.Uint64(block[16:24])
	// 	t = t + 0x10000
	// 	block, err = LoadBlock(t, GetCurrentTimestamp())
	// }
	// hash := cryptography.HashTwice(res)
	// return Hash2Address(hash), nil
}
