# Cross Platform XDAG GUI Wallet

This is a cross-platform XDAG GUI wallet, especially for macOS and Linux, powered by [fyne](https://github.com/fyne-io/fyne).

Fyne is a cross-platform GUI in Go inspired by Material Design.

The wallet can run on Windows, Linux, Mac now.

Mobile version is coming.
## usage
- usage: by command-line parameter -help
- 3 run modes: gui(default), cli, server, by commnad-line parameter -mode
- cli commands:
  - help -- display commands list
  - exit -- exit cli wallet
  - account -- display address of wallet account
  - balance -- display balance of wallet account
  - xfer V A R -- transfer V coins to address A with remark R
  - mnemonic -- display mnemonic of wallet account
  - export P -- export mnemonic to file P
- jsonrpc server: by command-line parameter -mode server -ip \<ip address\> -port \<port number\>
  - method: Xdag.Unlock
    - params: ["\<wallet password\>"]
    - response: {"id":1,"result":"success","error":null}
  - method: Xdag.Lock
    - params: ["\<wallet password\>"]
    - response: {"id":1,"result":"success","error":null}
  - method: Xdag.Account
    - params: [""]
    - response: {"id":1,"result": "\<wallet address\>","error":null}
  - method: Xdag.Balance
    - params: [""]
    - response: {"id":1,"result": "\<wallet balance\>","error":null}
  - method: Xdag.Balance
    - params: ["\<wallet address\>"]
    - response: {"id":1,"result": "\<balance of the address\>","error":null}
  - method: Xdag.Transfer
    - params: [{"amount":"\<amount\>","address":"\<to address\>","remark":"\<remark\>"}]
    - response: {"id":1,"result": {"Status":"success","TxHash":"\<transaction hash\>"},"error":null}

## repo structure
 - clib - a wrapper of XDAG Wallet C library
   - xDagWallet - XDAG wallet C library
 - wallet - golang XDAG wallet app
   - i18n - international strings
   - data - i18n config json, fonts
   - images - image and icon bundled in components/resource.go
   - component - ui of wallet window
   - config - wallet config
   - wallet_state - wallet state
   - xlog - wallet log
   - xdago - bip32,bip39,bip44
   - cli - command line mode
   - server - rpc server mode

## build
enter /clib

build runtime library with CMakeLists.txt

need MingW64 in Windows

enter /wallet

`$ go mod tidy`

`$ CGO_ENABLED=1 go build`

in Windows

`> CGO_ENABLED=1 go build -ldflags -H=windowsgui`

in Mac

if library not found , go to /wallet/components/wallet_cgo.go

change the library path in code:

`//#cgo darwin LDFLAGS: ....`

## deployment
enter /wallet

copy goXdagWallet(.exe), wallet-config.json and data folder to your deployment path.

### Windows

also need copy libcrypto-1_1-x64.dll and libwinpthread-1.dll in MingW64's bin path to deployment path.

### Linux and Mac

need install secp256k1 and openssl first

- Ubuntu and Debian:

`$ sudo apt-get install libsecp256k1-dev openssl libssl-dev`

- Fedora and Centos:

`$ sudo yum install openssl openssl-devel`

download and build from source

`github.com/bitcoin-core/secp256k1.git`


- Manjaro and Arch linux:

`$ sudo pacman libsecp256k1 openssl`

- Mac

`$ brew install openssl`

`$ echo 'export PATH="/usr/local/opt/openssl/bin:$PATH"' >> ~/.bash_profile`

`$ source ~/.bash_profile`

build secp256k1 from source

`github.com/bitcoin-core/secp256k1.git`

## features

 - frequently transferred addresses list
 - history query filter and pagination
