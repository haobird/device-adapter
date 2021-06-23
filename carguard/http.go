package carguard

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/haobird/goutils"
	"github.com/haobird/logger"
	"github.com/tidwall/gjson"
)

// HTTP 接口相关
// 接口请求参数
type ReqParams struct {
	Topic   string
	Key     string
	Payload interface{}
}

type responseInfo struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

func respondWithInfo(code int, message string, data interface{}, c *gin.Context) {
	c.JSON(200, responseInfo{
		Code: code,
		Msg:  message,
		Data: data,
	})
	c.Abort()
}

//InitHTTP 开放对外接口
func InitHTTP() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// 基础数据上报
	router.POST("/api/upark/basicinfo", func(c *gin.Context) {
		// 解析 json 数据
		buf, err := c.GetRawData()
		// 处理业务
		if err != nil {
			respondWithInfo(101, err.Error(), nil, c)
			return
		}

		go ProcessPublsihRaw("basicinfo", buf)
		respondWithInfo(200, "success", nil, c)
	})

	// 心跳保活上报
	router.POST("/api/upark/keepalive", func(c *gin.Context) {
		// 解析 json 数据
		buf, err := c.GetRawData()
		// 处理业务
		if err != nil {
			respondWithInfo(101, err.Error(), nil, c)
			return
		}

		go ProcessPublsihRaw("keepalive", buf)
		respondWithInfo(200, "success", nil, c)
	})

	// 车位告警上报
	router.POST("/api/upark/parkalarm", func(c *gin.Context) {
		// 解析 json 数据
		buf, err := c.GetRawData()
		// 处理业务
		if err != nil {
			respondWithInfo(101, err.Error(), nil, c)
			return
		}

		go ProcessPublsihRaw("parkalarm", buf)
		respondWithInfo(200, "success", nil, c)
	})

	// 过车抓拍上报
	router.POST("/api/upark/capture", func(c *gin.Context) {
		// 解析 json 数据
		buf, err := c.GetRawData()
		// 处理业务
		if err != nil {
			respondWithInfo(101, err.Error(), nil, c)
			return
		}

		go ProcessPublsihRaw("capture", buf)
		respondWithInfo(200, "success", nil, c)
	})

	// 开闸放行结果上报
	router.POST("/api/upark/notifyresult/gatecontrol", func(c *gin.Context) {
		// 解析 json 数据
		buf, err := c.GetRawData()
		// 处理业务
		if err != nil {
			respondWithInfo(101, err.Error(), nil, c)
			return
		}

		go ProcessPubackRaw("gatecontrol", buf)
		respondWithInfo(200, "success", nil, c)
	})

	// 手动抓拍结果上报
	router.POST("/api/upark/notifyresult/manualcapture/common", func(c *gin.Context) {
		// 解析 json 数据
		buf, err := c.GetRawData()
		// 处理业务
		if err != nil {
			respondWithInfo(101, err.Error(), nil, c)
			return
		}

		go ProcessPubackRaw("manualcapturecommon", buf)
		respondWithInfo(200, "success", nil, c)
	})

	// 车位灯设置结果上报
	router.POST("/api/upark/notifyresult/lamp", func(c *gin.Context) {
		// 解析 json 数据
		buf, err := c.GetRawData()
		// 处理业务
		if err != nil {
			respondWithInfo(101, err.Error(), nil, c)
			return
		}

		go ProcessPubackRaw("lamp", buf)
		respondWithInfo(200, "success", nil, c)
	})

	// 黑白名单同步结果上报
	router.POST("/api/upark/notifyresult/list", func(c *gin.Context) {
		// 解析 json 数据
		buf, err := c.GetRawData()
		// 处理业务
		if err != nil {
			respondWithInfo(101, err.Error(), nil, c)
			return
		}

		go ProcessPubackRaw("authorized", buf)
		respondWithInfo(200, "success", nil, c)
	})

	// 开闸放行结果上报
	router.POST("/api/upark/carguard", func(c *gin.Context) {
		// 解析 json 数据
		buf, err := c.GetRawData()
		// 处理业务
		if err != nil {
			respondWithInfo(101, err.Error(), nil, c)
			return
		}

		code, content, err := AsyncReq(buf)

		statusCode := gjson.Get(content, "Response.StatusCode").Int()
		if statusCode != 0 {
			statusString := gjson.Get(content, "Response.StatusString").String()
			respondWithInfo(int(statusCode), statusString, content, c)
			return
		}

		responseCode := gjson.Get(content, "Response.ResponseCode").Int()
		if responseCode != 0 {
			responseString := gjson.Get(content, "Response.ResponseString").String()
			respondWithInfo(int(responseCode), responseString, content, c)
			return
		}

		// respondWithInfo(0, "success", string(content), c)
		respondWithInfo(code, err.Error(), content, c)
	})

	httpAddr = config.HTTPAddr
	logger.Debug("http服务启动: ", httpAddr)
	router.Run(httpAddr)
}

//AsyncReq 处理异步回调请求: 获取 client, 发送封包，等待响应，返回结果
func AsyncReq(buf []byte) (int, string, error) {
	var err error
	var input ReqParams
	err = json.Unmarshal(buf, &input)
	if err != nil {
		return 1, "", errors.New("参数解析错误")
	}

	fmt.Println("api接收到", string(buf))
	data, err := json.Marshal(input.Payload)
	// 生成消息 id
	msgId := goutils.String(time.Now().Unix())
	logger.Debug("msgID :", msgId)

	// 封装Package
	pack := Package{
		MessageType: Command,
		Topic:       input.Topic,
		RequestID:   msgId,
		ClientID:    input.Key,
		Data:        data,
	}

	if err != nil {
		return 2, "", err
	}

	// 建立消息通道
	ch := make(chan string, 1)
	msgChans[msgId] = ch

	// 执行命令下发
	go ProcessCommandPack(pack)

	// 循环读取 通道的响应，并增加超时退出
	var code int = 0
	var msg string
	select {
	case msg = <-ch:
		fmt.Println("data from channel ", msgId)
		err = errors.New("success")
	case <-time.After(10 * time.Second):
		fmt.Println("TimeOut")
		err = errors.New("响应超时TimeOut")
		code = 3
	}

	// 关闭通道，删除 key
	delete(msgChans, msgId)
	close(ch)

	// 返回结果
	return code, msg, err
}
