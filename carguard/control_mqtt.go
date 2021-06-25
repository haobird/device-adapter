package carguard

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/haobird/goutils"
	"github.com/haobird/logger"
)

type MQTTControl struct {
	conn mqtt.Client
}

func (m *MQTTControl) Publish(packet *Package) error {
	topic := packet.Topic
	data := packet.Data
	fmt.Println("topic:", topic)
	fmt.Println("data:", string(data))
	var input map[string]interface{}
	err := json.Unmarshal(data, &input)

	if err != nil {
		return err
	}

	// 判断 mqtt的模式
	if config.Mqtt["mode"] == "device" {
		return m.PublishDevice(input)
	} else {
		return m.PublishSpecial(input)
	}
}

// 发送 设备能够识别的 消息包
func (m *MQTTControl) PublishDevice(input map[string]interface{}) error {
	var qoss byte = 0
	buf, err := json.Marshal(input["payload"])
	fmt.Println("发送topic:", config.Mqtt["topicPrefix"], "发送 body:", string(buf))
	m.conn.Publish(config.Mqtt["topicPrefix"], qoss, false, buf)
	return err
}

// 发送 本地服务能够识别的 消息包
func (m *MQTTControl) PublishSpecial(input map[string]interface{}) error {
	var qoss byte = 0
	// 删除冗余的payload , 剩下的发送
	delete(input, "payload")
	buf, err := json.Marshal(input)
	topic := config.Mqtt["topicPrefix"] + "_special"
	m.conn.Publish(topic, qoss, false, buf)
	return err
}

func InitMQTTControl() *MQTTControl {
	var option = config.Mqtt
	// json.Unmarshal([]byte(config.Mqtt), &option)
	logger.Info("mqtt配置为：", option)
	clientid := option["clientid"]
	salt, _ := goutils.GenShortID()
	option["clientid"] = clientid + "_" + salt
	client, _ := NewMQTTClient(option)
	return &MQTTControl{
		conn: client,
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	logger.Debug("MQTT Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	logger.Debug("Connect lost: %v", err)
	RetryConnect(client)
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	logger.Debugf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	// 解析分割类型
	// topic := msg.Topic()
	// 格式化返回的数据

	// 调用相应的回调

}

func NewMQTTClient(option map[string]string) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(option["addr"])
	opts.SetClientID(option["clientID"])
	opts.SetUsername(option["username"])
	opts.SetPassword(option["password"])
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	conn := mqtt.NewClient(opts)
	// RetryConnect(conn)
	if token := conn.Connect(); token.Wait() && token.Error() != nil {
		// 打印错误日志
		return nil, token.Error()
	}
	return conn, nil
}

func RetryConnect(conn mqtt.Client) {
	for {

		if token := conn.Connect(); token.Wait() && token.Error() != nil {
			// 打印错误日志
			return
		}
		fmt.Println("[control mqtt] RetryConnect ")
		time.Sleep(3 * time.Second)
	}

}
