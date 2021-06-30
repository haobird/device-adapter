package guardface

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/haobird/logger"
)

//Config 配置结构体
type Config struct {
	TCPAddr  string   `json:"tcpAddr"`
	HTTPAddr string   `json:"httpAddr"`
	AMQP     AMQPConf `json:"amqp"`
	Bridge   string   `json:"bridge"`
	Log      LogConf
}

type LogConf struct {
	Writers       string `json:"writers"`        // file,stdout  # 文件和终端输出
	Level         string `json:"level"`          // debug    # 报警等级
	File          string `json:"file"`           // /data/log/lite.log
	FormatText    bool   `json:"format_text"`    // false
	Color         bool   `json:"color"`          // false
	RollingPolicy string `json:"rolling_policy"` // size
	RotateDate    int    `json:"rotate_date"`
	RotateSize    int    `json:"rotate_size"`
	BackupCount   int    `json:"backup_count"`
}

//AMQPConf rabbitmq配置
type AMQPConf struct {
	Addr           string        `json:"addr"`
	PublishChannel QueueExchange `json:"onPublish"`
}

//QueueExchange 交换机结构体
type QueueExchange struct {
	QueueName    string `json:"queue_name"`    // 队列名称
	RoutingKey   string `json:"routing_key"`   // key值
	ExchangeName string `json:"exchange_name"` // 交换机名称
	ExchangeType string `json:"exchange_type"` // 交换机类型
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
