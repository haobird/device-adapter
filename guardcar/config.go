package guardcar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/haobird/logger"
)

//Config 配置结构体
type Config struct {
	TCPAddr     string            `json:"tcpAddr"`
	HTTPAddr    string            `json:"httpAddr"`
	Bridge      string            `json:"bridge"`
	ControlMode string            `json:"controlMode"`
	Mqtt        map[string]string `json:"mqtt"`
}

func LoadConfig(path string) *Config {
	logger.Info("start load config....")
	content, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Fatal("Read config file error: ", err)
	}
	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		logger.Fatal("Unmarshal config file error: ", err)
	}
	fmt.Println(config)

	return &config
}
