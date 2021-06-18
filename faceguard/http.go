package faceguard

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
	Msg  string      `json:"msg"`
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

	router.POST("api/faceguard", func(c *gin.Context) {
		// 解析 json 数据
		buf, _ := c.GetRawData()
		code, content, err := AsyncReq(buf)
		if err != nil {
			respondWithInfo(code, err.Error(), content, c)
			return
		}

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

		respondWithInfo(0, "success", string(content), c)

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
	// 读取client是否存在
	cli := manager.GetClient(input.Key)
	if cli == nil {
		return 2, "", errors.New("设备不在线")
	}
	// RespAuthorized(cli)
	// return 2, "", err
	// 经过sdk转换为设备可以识别的请求
	buf, err = json.Marshal(input.Payload)
	// fmt.Println("api接收到", string(buf))
	// 生成消息 id
	msgId := goutils.String(time.Now().Unix())
	logger.Debug("msgID :", msgId)
	// 封装Package
	pack := Package{
		Topic:     input.Topic,
		RequestID: msgId,
		ClientID:  input.Key,
		Data:      buf,
	}
	content := sdk.Trans(pack)
	// fmt.Println(string(content))
	// fmt.Println("字符串长度", len(string(content)))
	err = cli.Write([]byte(content))

	if err != nil {
		return 2, "", err
	}

	// 建立消息通道
	ch := make(chan string, 1)
	msgChans[msgId] = ch

	// 循环读取 通道的响应，并增加超时退出
	var code int = 0
	var msg string
	select {
	case msg = <-ch:
		fmt.Println("data from channel ", msgId)
	// 如果把这个注释掉，则会阻塞 deadlock
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
