package conf

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/toolkits/file"
)

type GlobalConfig struct {
	Debug    bool           `json:"debug"`
	Database string         `json:"database"`
	Http     *HttpConfig    `json:http`
	Redis    *RedisConfig   `json:redis`
	JsonRpc  *JsonRpcConfig `json:jsonprc`
}

type HttpConfig struct {
	Port string `json:"port"`
}

type RedisConfig struct {
	Addr string `json:"addr"`
}

type JsonRpcConfig struct {
	Addrs []string `json:"addrs"`
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
