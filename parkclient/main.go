package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/haobird/logger"
	"github.com/tidwall/gjson"
)

// 配合 云端 适配器使用的 本地适配器

var (
	ch      = make(chan int)
	cfgFile = flag.String("cfg", "config.json", "config file")
)

func main() {
	flag.Parse()
	// 读取配置文件
	logger.Info("start load config....", *cfgFile)
	content, err := ioutil.ReadFile(*cfgFile)
	if err != nil {
		logger.Fatal("Read config file error: ", err)
	}
	var option map[string]string
	err = json.Unmarshal(content, &option)
	if err != nil {
		logger.Fatal("Unmarshal config file error: ", err)
	}

	// 重复建立 mqtt连接
	go func() {
		for {
			select {
			case <-ch:
				client, err := NewMQTTClient(option)
				if err != nil {
					ch <- 1
				}
				// 订阅 topic
				Subscribe(client, option["topicPrefix"], 0, messagePubHandler)
			}
		}
	}()
	ch <- 1

	keepAlive()

}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	logger.Debug("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	logger.Debug("Connect lost: %v", err)
	ch <- 1
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	logger.Debugf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	// 解析分割类型
	// topic := msg.Topic()
	// 格式化返回的数据
	// var input map[string]interface{}
	// err := json.Unmarshal(msg.Payload(), &input)
	// if err != nil {
	// 	logger.Error("订阅获取到的消息错误:", err)
	// }

	// buf,_ := json.Marshal(input["body"])
	url := gjson.GetBytes(msg.Payload(), "url").String()
	method := gjson.GetBytes(msg.Payload(), "method").String()
	body := gjson.GetBytes(msg.Payload(), "body").String()

	// 调用相应的回调
	resp, err := Request(url, method, []byte(body), nil)
	logger.Debug("请求设备响应结果：", resp, ", 错误: ", err)
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
	if token := conn.Connect(); token.Wait() && token.Error() != nil {
		// 打印错误日志
		// panic(token.Error())
		return nil, token.Error()
	}
	return conn, nil
}

func Subscribe(conn mqtt.Client, topic string, qos byte, callback mqtt.MessageHandler) error {
	conn.Subscribe(topic, qos, callback)
	return nil
}

//Request 发起 Http请求
func Request(url string, method string, data []byte, headers map[string]string) (result string, err error) {
	client := &http.Client{Timeout: 10 * time.Second}

	// jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodystr := string(body)
	return bodystr, nil

	// if resp.StatusCode == 200 {
	// 	body, err_ := ioutil.ReadAll(resp.Body)
	// 	if err_ != nil {
	// 		return "", err_
	// 	}
	// 	bodystr := string(body)
	// 	return bodystr, nil
	// }
	// return "", err

}

func keepAlive() {
	//合建chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	//阻塞直到有信号传入
	fmt.Println("总进程服务启动完成")
	//阻塞直至有信号传入
	s := <-c
	fmt.Println("退出信号", s)
}
