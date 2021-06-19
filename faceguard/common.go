package faceguard

import (
	"deviceadapter/bridge"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/haobird/fixpool"
	"github.com/haobird/logger"
)

// 数据包的几种状态
const (
	Unknown    = "unknown"
	Connect    = "connect"
	Heart      = "heart"
	Publish    = "publish"
	PubAck     = "puback"
	Command    = "command"
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
	lifeSpan = 60              // 单位 秒
	wpool    = fixpool.New(20) // 设备任务池
	tpool    = fixpool.New(20) // 请求任务池
)

//Message 消息
type Message struct {
	client *Client
	packet *Package
}

func Init(cfgFile string) {
	// 加载 配置文件
	config = LoadConfig(cfgFile)

	// 创建桥接
	mybridge = bridge.NewBridge(config.Bridge)

	// 建立 tcp连接
	go InitTCP()

	// 建立 API 接口
	go InitHTTP()

	// 保持进程
	keepAlive()
}

// ProcessCommand 处理命令
func ProcessCommand(c *Client, p []byte) {
	wpool.Submit(c.ID, func() {
		c.Write(p)
	})
}

// ProcessDeviceData 处理数据
func ProcessDeviceData(c *Client, p []byte) {
	logger.Debugf("client %s ProcessMessage %s", c.ID, string(p))
	// 解析数据包
	packet, err := sdk.HandlePacket(p)
	if err != nil {
		logger.Errorf("client %s ProcessMessage error:%s", c.ID, err)
		return
	}
	msg := Message{
		client: c,
		packet: packet,
	}
	ProcessMessage(msg)
}

// 处理消息
func ProcessMessage(msg Message) {
	c := msg.client
	packet := msg.packet

	// 基于消息内容 进行 相应的处理
	messageType := packet.MessageType
	switch messageType {
	case Connect:
		logger.Infof("[%s] 设备Register ", c.ID)
		RegisterHandler(c)
	case PubAck:
		PubAckHandler(packet)
	case Publish:
		PublishHandler(packet)
	case Heart:
		HeartHandler(c)
	case Disconnect:
		DisconnectHandler(c)
	}
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
