package guardface

import (
	"deviceadapter/bridge"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/haobird/fixpool"
	"github.com/haobird/logger"
)

var (
	config   *Config
	manager  = &ConnManager{}
	tcpAddr  = ":3570"
	httpAddr = ":9081"
	sdk      = &Things{}
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

// ProcessDeviceData 处理数据
func ProcessRawData(c *Client, p []byte) {
	// str := string(p)
	// 过滤掉图片数据
	// p = []byte(Tidy(str))
	logger.Debugf("client %s ProcessRawData %s", c.ID, string(p))

	// 解析数据包
	packet := sdk.ProcessDataUp(p)
	if packet == nil {
		logger.Errorf("client %s ProcessMessage error", c.ID)
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
	messageType := packet.MessageType

	logger.Debugf("client %s ProcessMessage %s", c.ID, messageType)

	// 基于消息内容 进行 相应的处理
	packet.ClientID = c.ID
	switch messageType {
	// case Connect:
	// 	logger.Infof("[%s] 设备Register ", c.ID)
	// 	RegisterHandler(c)
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

// 发送到设备 数据包
func SubmitWork(c *Client, packet *Package) {
	logger.Debugf("SubmitWork设备数据：MessageType:%s,RequestID:%s,ClientID:%s,Topic:%s, Data:%s",
		packet.MessageType,
		packet.RequestID,
		packet.ClientID,
		packet.Topic,
		string(packet.Data))
	// buf, _ := json.Marshal(packet)
	ProcessCommand(c, packet.Data)
}

// 执行下发任务
func ProcessCommand(c *Client, p []byte) {
	wpool.Submit(c.ID, func() {
		c.Write(p)
	})
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
