package carguard

import (
	"deviceadapter/bridge"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/haobird/fixpool"
	"github.com/haobird/goutils"
)

var (
	config   *Config
	httpAddr = ":9081"
	mybridge bridge.Bridge
	wpool    = fixpool.New(20) // 设备任务池
	things   = &Things{}
	msgChans = make(map[string]chan string)
	control  Control
)

// 服务启动
func Init(cfgFile string) {
	// 加载 配置文件
	config = LoadConfig(cfgFile)

	// 创建桥接
	mybridge = bridge.NewBridge(config.Bridge)

	control = NewControl(config.ControlMode)

	// 建立 API 接口
	go InitHTTP()

	// 保持进程
	keepAlive()
}

// ProcessPublsih 处理上行消息
func ProcessPublsihRaw(action string, buf []byte) {
	packet := things.ParsePublishData(action, buf)
	messageType := packet.MessageType
	if messageType == Heart {
		// 判断是否存在当前设备的缓存
		// 不存在，则执行注册 方法， 复写当前的包
		packet = things.ParsePublishData(Connect, buf)
		messageType = packet.MessageType

	}

	if messageType == Publish {
		// 执行上报逻辑
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

}

// ProcessPubackRaw 处理响应消息
func ProcessPubackRaw(action string, buf []byte) {
	packet := things.ParsePubackData(action, buf)
	requestID := packet.RequestID
	if ch, ok := msgChans[requestID]; ok {
		ch <- string(packet.Data)
	}
}

// 处理指令下发
func ProcessCommandPack(packet Package) {
	// 格式化数据
	newPacket := things.ParseCommanData(packet)
	// 执行下发
	control.Publish(newPacket)
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
