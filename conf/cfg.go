package conf

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/MrYang/golang-learn/kit/file"
)

type GlobalConfig struct {
	Common *CommonConfig `json:"common"`
	Server *ServerConfig `json:"server"`
	Client *ClientConfig `json:"client"`
}

type CommonConfig struct {
	Redis    string `json:"redis"`
	Database string `json:"database"`
}

type ServerConfig struct {
	Http    *ServerHttpConfig    `json:"http"`
	JsonRpc *ServerJsonRpcConfig `json:"jsonRpc"`
	Tcp     *TcpConfig           `json:"tcp"`
	Rpc     *RpcConfig           `json:"rpc"`
	GRpc    *GRpcConfig          `json:"gRpc"`
}

type ServerHttpConfig struct {
	Port string `json:"port"`
}

type ServerJsonRpcConfig struct {
	Listen string `json:"listen"`
}

type RpcConfig struct {
	Listen string `json:"listen"`
}

type TcpConfig struct {
	Listen string `json:"listen"`
}

type GRpcConfig struct {
	Listen string `json:"listen"`
}

type ClientConfig struct {
	JsonRpc string `json:"jsonRpc"`
	GRpc    string `json:"gRpc"`
	Rpc     string `json:"rpc"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.Lock()
	defer configLock.Unlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not exist")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)

	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &c
	log.Println("read config file:", cfg, "successfully")
}
