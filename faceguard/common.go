package faceguard

import (
	"deviceadapter/bridge"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/haobird/fixpool"
	"github.com/haobird/goutils"
)

// 连接内容的几种状态
const (
	Unknown    = "unknown"
	Connect    = "connect"
	Heart      = "heart"
	Publish    = "publish"
	PubAck     = "puback"
	Disconnect = "disconnect"
)

var (
	config   *Config
	manager  = &ConnManager{}
	tcpAddr  = ":3570"
	httpAddr = ":9081"
	sdk      = &Sdk{}
	msgChans = make(map[string]chan string)
	mybridge bridge.Bridge
	lifeSpan = 60 // 单位 秒
	wpool    = fixpool.New(20)
)

func Init() {
	// 加载 配置文件
	config = LoadConfig("config.json")

	// 创建桥接
	mybridge = bridge.NewBridge("")

	// 建立 tcp连接
	go InitTCP()

	// 建立 API 接口
	go InitHTTP()

	// 保持进程
	keepAlive()
}

// BusinessHandler 处理 业务逻辑
func BusinessHandler(packet *Package) {
	if packet.MessageType == PubAck {
		requestID := packet.RequestID
		if ch, ok := msgChans[requestID]; ok {
			ch <- string(packet.Data)
		}
	}
	// fmt.Println(body)
	ele := &bridge.Element{
		MessageType: packet.MessageType,
		RequestID:   packet.RequestID,
		Timestamp:   goutils.Int64(packet.RequestID),
		ClientID:    packet.ClientID,
		Topic:       packet.Topic,
		Data:        string(packet.Data),
	}
	mybridge.Publish(ele)
}

// 处理 客户端 关闭逻辑
func BeforeCloseHandler(clientid string) {
	manager.DeleteClient(clientid)
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
