package utils

import (
	"time"
)

const (
	XDAG_TEST_ERA = 0x16900000000
	XDAG_MAIN_ERA = 0x16940000000
)

// GetEndOfEpoch 获取时间戳所属epoch的最后一个时间戳 主要用于mainblock
func GetEndOfEpoch(t uint64) uint64 {
	return t | 0xffff
}

func IsEndOfEpoch(t uint64) bool {
	return (t & 0xffff) == 0xffff
}

// GetEpoch 获取该时间戳所属的epoch
func GetEpoch(t uint64) uint64 {
	return t >> 16
}

// GetCurrentTimestamp 获取当前的xdag时间戳
func GetCurrentTimestamp() uint64 {
	t := time.Now().UTC().UnixNano()
	sec := t / 1e9
	usec := (t - sec*1e9) / 1e3
	xmsec := (usec << 10) / 1e6
	return uint64(sec)<<10 | uint64(xmsec)
}

// Ms2XdagTimestamp 把毫秒转为xdag时间戳
func Ms2XdagTimestamp(ms uint64) uint64 {
	sec := ms / 1e3
	xmsec := ((ms - sec*1e3) << 10) / 1e3
	return (sec << 10) | xmsec
}

func XdagTimestamp2Ms(t uint64) uint64 {
	sec := t >> 10
	xms := t - (sec << 10)
	ms := (xms * 1e3) >> 10
	return (sec * 1e3) + ms
}

// GetMainTime 获取当前时间所属epoch的最后一个时间戳
func GetMainTime() uint64 {
	return GetEndOfEpoch(GetCurrentTimestamp())
}

func GetCurrentEpoch() uint64 {
	return GetEpoch(GetCurrentTimestamp())
}
