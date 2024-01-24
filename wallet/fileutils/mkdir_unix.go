//go:build !windows

package fileutils

import (
	"os"

	"golang.org/x/sys/unix"
)

func MkdirAll(p string) error {
	mask := unix.Umask(0)  // umask 0000
	defer unix.Umask(mask) // recover umask
	return os.MkdirAll(p, 0777)
}

func WriteFile(p string, data []byte) error {
	mask := unix.Umask(0)  // umask 0000
	defer unix.Umask(mask) // recover umask
	return os.WriteFile(p, data, 0777)
}
