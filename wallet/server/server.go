package server

import (
	"errors"
	"fmt"
	"goXdagWallet/components"
	"goXdagWallet/xdago/base58"
	"goXdagWallet/xdago/cryptography"
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"sync"
)

var HelloServiceName = "Xdag"

type XferParam struct {
	Amount  string `json:"amount"`
	Address string `json:"address"`
	Remark  string `json:"remark"`
}

type XferRespone struct {
	Status string `json:"Status"`
	TxHash string `json:"TxHash"`
}

type Xdag struct {
	lock sync.Mutex
}

func (s *Xdag) Unlock(password string, reply *string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	res := components.ConnectBipWallet(password)
	if res {
		components.PwdStr = password
		b := cryptography.ToBytesAddress(components.BipWallet.GetDefKey())
		components.BipAddress = base58.ChkEnc(b[:])
	} else {
		return errors.New("incorrect password")
	}
	*reply = "success"
	return nil
}

func (s *Xdag) Lock(password string, reply *string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if components.PwdStr == "" || components.BipAddress == "" {
		*reply = "wallet locked"
		return nil
	}
	if components.BipWallet == nil {
		*reply = "wallet locked"
		return nil
	}
	if password != components.PwdStr {
		return errors.New("incorrect password")
	}

	components.BipWallet.LockWallet()
	components.PwdStr = ""
	components.BipAddress = ""
	*reply = "success"
	return nil
}

func (s *Xdag) Account(request string, reply *string) error {
	if components.PwdStr == "" || components.BipAddress == "" ||
		components.BipWallet == nil || components.BipWallet.IsLocked() {
		return errors.New("wallet is locked")
	}
	*reply = components.BipAddress
	return nil
}

func (s *Xdag) Transfer(request XferParam, reply *XferRespone) error {
	if components.PwdStr == "" || components.BipAddress == "" ||
		components.BipWallet == nil || components.BipWallet.IsLocked() {
		return errors.New("wallet is locked")
	}
	if !components.ValidateBipAddress(request.Address) {
		return errors.New("address format error")
	}
	value, err := strconv.ParseFloat(request.Amount, 64)
	if err != nil || value <= 0.0 {
		return errors.New("amount number error")
	}
	if !components.ValidateRemark(request.Remark) {
		return errors.New("remark format error")
	}

	fromValue, err := components.BalanceRpc(components.BipAddress)
	if err != nil {
		return err
	}
	balance, _ := strconv.ParseFloat(fromValue, 64)
	if balance < value {
		return errors.New("insufficient amount")
	}
	hash, err := components.TransferRpc(components.BipAddress, request.Address,
		request.Amount, request.Remark, components.BipWallet.GetDefKey())
	if err != nil {
		return err
	}
	// *reply = "success:" + hash
	*reply = XferRespone{
		Status: "success",
		TxHash: hash,
	}
	return nil
}

func (s *Xdag) Balance(address string, reply *string) error {
	if components.PwdStr == "" || components.BipAddress == "" ||
		components.BipWallet == nil || components.BipWallet.IsLocked() {
		return errors.New("wallet is locked")
	}

	if address == "" {
		address = components.BipAddress
	} else if !components.ValidateBipAddress(address) {
		return errors.New("invalid address")
	}
	balance, err := components.BalanceRpc(address)
	if err != nil {
		return err
	}
	*reply = balance
	return nil
}

func RunServer(ip, port string) {
	rpc.RegisterName(HelloServiceName, new(Xdag))

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			ReadCloser: r.Body,
			Writer:     rw,
		}
		rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
	})
	fmt.Println("Listen and serve on", ip+":"+port, "...")
	http.ListenAndServe(ip+":"+port, nil)
}
