package config

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

const configFile = "wallet-config.json"

var conf Config

type Config struct {
	Option      WalletOption `json:"wallet_option"`
	Version     string       `json:"version"`
	CultureInfo string       `json:"culture_info"`
	Addresses   []string     `json:"addresses"`
	Query       DefaultQuery `json:"query"`
}

type WalletOption struct {
	RpcEnabled           bool    `json:"rpc_enabled"`
	RpcPort              *int    `json:"rpc_port"`
	SpecifiedPoolAddress *string `json:"specified_pool_address"`
	IsTestNet            bool    `json:"is_test_net"`
	DisableMining        bool    `json:"disable_mining"`
	PoolAddress          string  `json:"pool_address"`
	TestnetApiUrl        string  `json:"testnet_api_url"`
}
type DefaultQuery struct {
	AmountFrom string `json:"amount_from"`
	AmountTo   string `json:"amount_to"`
	Timestamp  string `json:"timestamp"`
	Remark     string `json:"remark"`
	Direction  string `json:"direction"`
}

func InitConfig() {
	newConf := false

	conf.Version = "0.5.5"
	conf.CultureInfo = "en-US"
	conf.Option.DisableMining = true
	conf.Option.PoolAddress = "xdag.org:13656"

	pwd, _ := os.Executable()
	pwd, _ = path.Split(pwd)
	data, err := os.ReadFile(path.Join(pwd, configFile))
	if err == nil {
		err = json.Unmarshal(data, &conf)
		if err != nil {
			log.Println("fail to Unmarshal configure,conf.json")
			newConf = true
		}
	} else {
		newConf = true
	}
	if conf.CultureInfo == "" {
		conf.CultureInfo = "en-US"
	}
	if newConf {
		data, _ = json.MarshalIndent(conf, "", "  ")
		os.WriteFile(configFile, data, 666)
	}
}
func DeleteAddress(id int) {
	if id >= 0 && id < len(conf.Addresses) {
		conf.Addresses = append(conf.Addresses[:id], conf.Addresses[id+1:]...)
		SaveConfig()
	}
}

func InsertAddress(address string) {
	pos := -1
	for i, item := range conf.Addresses {
		if item == address {
			pos = i
			break
		}
	}
	if pos > 0 {
		conf.Addresses = append(conf.Addresses[:pos], conf.Addresses[pos+1:]...)
		conf.Addresses = append([]string{address}, conf.Addresses...)
	} else if pos == -1 {
		conf.Addresses = append([]string{address}, conf.Addresses...)
		if len(conf.Addresses) > 10 {
			conf.Addresses = conf.Addresses[:10]
		}
	}
	if pos != 0 {
		SaveConfig()
	}

}

func SaveConfig() error {
	data, _ := json.MarshalIndent(conf, "", "  ")
	os.WriteFile(configFile, data, 666)
	return nil
}

func GetConfig() *Config {
	return &conf
}
