//go:build windows

package fileutils

import "os"

func MkdirAll(p string) error {
	return os.MkdirAll(p, 0666)
}

func WriteFile(p string, data []byte) error {
	return os.WriteFile(p, data, 0666)
}
