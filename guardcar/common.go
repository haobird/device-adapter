package guardcar

import (
	"deviceadapter/bridge"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/haobird/fixpool"
	"github.com/haobird/goutils"
	"github.com/haobird/logger"
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

	// 格式化日志输出
	logger.InitWithConfig(&logger.Cfg{
		Writers:       config.Log.Writers,
		Level:         config.Log.Level,
		File:          config.Log.File,
		FormatText:    config.Log.FormatText,
		Color:         config.Log.Color,
		RollingPolicy: config.Log.RollingPolicy,
		RotateDate:    config.Log.RotateDate,
		RotateSize:    config.Log.RotateSize,
		BackupCount:   config.Log.BackupCount,
	})

	// 创建桥接
	mybridge = bridge.NewBridge(config.Bridge)

	control = NewControl(config.ControlMode)

	// 初始化缓存和回调函数
	cache = InitCache(func(item *CacheItem) {
		DisconnectHandler(item.key)
	})

	// 建立 API 接口
	go InitHTTP()

	// 保持进程
	keepAlive()
}

// ProcessPublsih 处理上行消息
func ProcessPublsihRaw(action string, buf []byte) {
	// str := string(buf)
	// 过滤掉图片数据
	// buf = []byte(Tidy(str))
	packet := things.ParsePublishData(action, buf)
	messageType := packet.MessageType
	clientID := packet.ClientID
	topic := packet.Topic

	if messageType == Heart {
		// 触发缓存更新
		cache.Value(clientID)
		return
	}

	// 如果是设备信息上报，则首先判断是否有缓存
	if topic == "deviceInfoCar" {
		// 判断是否存在当前设备的缓存
		_, err := cache.Value(clientID)
		if err != nil {
			// 不存在，则添加缓存
			cache.Add(clientID, 2*time.Minute, packet.Data)
			// 基础数据上报
			packet = things.ParsePublishData(Connect, buf)
			messageType = Publish
		} else {
			// 如果存在，则当做心跳数据
			// messageType = Heart
		}
	}

	packet.ClientID = clientID

	if messageType == Publish {
		// 执行上报逻辑
		ele := &bridge.Element{
			MessageType: "vehicledetection",
			RequestID:   packet.RequestID,
			Timestamp:   goutils.Int64(packet.RequestID),
			ClientID:    packet.ClientID,
			Topic:       packet.Topic,
			Data:        string(packet.Data),
		}
		mybridge.Publish(ele)
	}

}

// 处理http的异步上报
func HandlerNotifyresult(action string, c *gin.Context) {
	// 解析 json 数据
	buf, err := c.GetRawData()
	logger.Debugf("Notifyresult [%s], body: %s", action, string(buf))
	// 处理业务
	if err != nil {
		respondWithInfo(101, err.Error(), nil, c)
		return
	}

	go ProcessPubackRaw(action, buf)
	respondWithInfo(0, "success", nil, c)

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
	logger.Debug("执行下发指令包的处理：", packet)
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
