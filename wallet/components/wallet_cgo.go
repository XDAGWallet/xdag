package components

//#cgo darwin LDFLAGS: -L${SRCDIR}/../../clib -lxdag_runtime_Darwin -L/usr/lib -lsecp256k1 -lm -L/usr/local/opt/openssl/lib -lssl -lcrypto
//#cgo linux LDFLAGS: -L${SRCDIR}/../../clib -lxdag_runtime_Linux -L/usr/lib -lsecp256k1 -lssl -lcrypto -lm
//#cgo windows LDFLAGS: -L${SRCDIR}/../../clib -lxdag_runtime_Windows -L/usr/lib -L/usr/local/lib -lsecp256k1 -lssl -lcrypto -lm -lws2_32
//#include "../../clib/xdag_runtime.h"
//#include "callback.h"
//#include <stdlib.h>
//#include <string.h>
/*
 typedef const char cchar_t;
*/
import "C"
import (
	"fmt"
	"goXdagWallet/config"
	"goXdagWallet/xdago/secp256k1"
	xdagoUtils "goXdagWallet/xdago/utils"
	"goXdagWallet/xlog"
	"unsafe"
)

func ConnectXdagWallet() int32 {
	var testnet int
	if config.GetConfig().Option.IsTestNet {
		testnet = 1
	}
	res := C.init_password_callback(C.int(testnet))
	result := int32(res)
	if result == 0 {
		xlog.Info("Initializing cryptography...")
		xlog.Info("Reading wallet...")
		k := getDefaultKey()
		if k == nil {
			xlog.Error("get default key failed.")
			fmt.Println("get default key failed.")
			return -4
		} else {
			XdagKey = secp256k1.PrivKeyFromBytes(k)
			addr, err := xdagoUtils.AddressFromStorage()
			if err != nil {
				xlog.Error(err)
				return -128
			} else if len(addr) < 1 {
				xlog.Error("wallet address not found")
				return -256
			} else if len(addr) == 1 {
				XdagAddress = addr[0]
			} else {
				xlog.Info(addr)
				OldAddresses = AddressWithBalance(addr)
				if len(OldAddresses) == 0 {
					xlog.Error("check old addresses balance failed")
					return -512
				} else if len(OldAddresses) == 1 {
					XdagAddress = OldAddresses[0]
				}
			}
		}
	}
	//fmt.Println(result)
	return result
}

//export goPasswordCallback
func goPasswordCallback(prompt *C.cchar_t, buf *C.char, size C.uint) C.int {
	C.memcpy(unsafe.Pointer(buf), unsafe.Pointer(&Password[0]), C.size_t(size))
	return C.int(0)
}

// get xdag wallet private key
func getDefaultKey() []byte {
	p := C.xdag_get_default_key()
	if uintptr(p) > 0 {
		key := C.GoBytes(p, 32)
		//fmt.Println(hex.EncodeToString(key[:]))
		//xlog.Info("default private key:", hex.EncodeToString(key[:]))
		return key
	}
	return nil
}
