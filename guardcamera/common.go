package guardcamera

import (
	"deviceadapter/bridge"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/haobird/goutils"
	"github.com/haobird/logger"
	"github.com/tidwall/gjson"
)

var (
	config   *Config
	httpAddr = ":9081"
	mybridge bridge.Bridge
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

	// 初始化缓存和回调函数
	cache = InitCache(func(item *CacheItem) {
		DisconnectHandler(item.key)
	})

	// 建立 API 接口
	go InitHTTP()

	// 保持进程
	keepAlive()
}

// 请求处理函数
func ReqHandler(action string, c *gin.Context) {
	buf, err := c.GetRawData()
	logger.Debugf("ReqHandler [%s], body: %s", action, string(buf))
	if err != nil {
		respondWithInfo(101, err.Error(), nil, c)
		return
	}
	var code = 0
	var msg = "success"
	var data string
	// 处理对应业务，并返回结果
	switch action {
	case "register":
		id := gjson.GetBytes(buf, "RegisterObject.DeviceID").String()
		curtime := goutils.GetTimeString(time.Now())
		data = fmt.Sprintf(`{"ResponseStatus":{"Id":"%s","LocalTime":"%s","RequestURL":"/VIID/System/Register","StatusCode":0,"StatusString":"OK"}}`, id, curtime)
		RegisterHandler(id)
	case "heart":
		id := gjson.GetBytes(buf, "KeepaliveObject.DeviceID").String()
		curtime := goutils.GetTimeString(time.Now())
		data = fmt.Sprintf(`{"ResponseStatus":{"Id":"%s","LocalTime":"%s","RequestURL":"/VIID/System/Keepalive","StatusCode":0,"StatusString":"OK"}}`, id, curtime)
		go HeartHandler(id)
	case "time":
		id := gjson.GetBytes(buf, "KeepaliveObject.DeviceID").String()
		curtime := goutils.GetTimeString(time.Now())
		data = fmt.Sprintf(`{"ResponseStatus":{"Id":"%s","LocalTime":"%s","RequestURL":"/VIID/System/Keepalive","StatusCode":0,"StatusString":"OK"}}`, id, curtime)
	case "faces":
		id := gjson.GetBytes(buf, "FaceListObject.FaceObject.0.DeviceID").String()
		curtime := goutils.GetTimeString(time.Now())
		data = fmt.Sprintf(`{"ResponseStatus":{"Id":"%s","LocalTime":"%s","RequestURL":"/VIID/Faces","StatusCode":0,"StatusString":"OK"}}`, id, curtime)
		go FacesHandler(id, buf)
	}
	var body map[string]interface{}
	json.Unmarshal([]byte(data), &body)
	respondWithInfo(code, msg, body, c)
}

// 处理消息上报
func ProcessPublsih(packet Package) {
	// 执行上报逻辑
	ele := &bridge.Element{
		MessageType: "cameraguard",
		RequestID:   packet.RequestID,
		Timestamp:   goutils.Int64(packet.RequestID),
		ClientID:    packet.ClientID,
		Topic:       packet.Topic,
		Data:        string(packet.Data),
	}
	mybridge.Publish(ele)
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
