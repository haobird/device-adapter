package guardcamera

import (
	"github.com/gin-gonic/gin"
	"github.com/haobird/logger"
)

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

	// 注册
	router.POST("/VIID/System/Register", func(c *gin.Context) { ReqHandler("register", c) })
	// 保活
	router.POST("/VIID/System/Keepalive", func(c *gin.Context) { ReqHandler("heart", c) })
	// 校时
	router.GET("/VIID/System/Time", func(c *gin.Context) { ReqHandler("time", c) })
	// 人脸抓拍
	router.POST("/VIID/Faces", func(c *gin.Context) { ReqHandler("faces", c) })

	// http服务启动
	httpAddr = config.HTTPAddr
	logger.Debug("http服务启动: ", httpAddr)
	router.Run(httpAddr)
}
