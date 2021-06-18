package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"deviceadapter/bridge"

	"github.com/gin-gonic/gin"
	utils "github.com/haobird/goutils"
	"github.com/haobird/logger"
)

const (
	ACCEPT_MIN_SLEEP = 100 * time.Millisecond
	ACCEPT_MAX_SLEEP = 10 * time.Second
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
)

//Config 配置结构体
type Config struct {
	TCPAddr  string   `json:"tcpAddr"`
	HTTPAddr string   `json:"httpAddr"`
	AMQP     AMQPConf `json:"amqp"`
}

//AMQPConf rabbitmq配置
type AMQPConf struct {
	Addr           string        `json:"addr"`
	PublishChannel QueueExchange `json:"onPublish"`
}

//QueueExchange 交换机结构体
type QueueExchange struct {
	QueueName    string `json:"queue_name"`    // 队列名称
	RoutingKey   string `json:"routing_key"`   // key值
	ExchangeName string `json:"exchange_name"` // 交换机名称
	ExchangeType string `json:"exchange_type"` // 交换机类型
}

func main() {
	// 加载 配置文件
	config = LoadConfig("unv/config.json")

	// 创建桥接
	mybridge = bridge.NewBridge("")

	// 建立 tcp连接
	go InitTCP()

	// 建立 API 接口
	go InitHTTP()

	// 保持进程
	keepAlive()
}

func LoadConfig(path string) *Config {
	logger.Info("start load config....")
	content, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Fatal("Read config file error: ", err)
	}
	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		logger.Fatal("Unmarshal config file error: ", err)
	}
	fmt.Println(config)

	return &config
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

// BusinessHandler 处理 业务逻辑
func BusinessHandler(ele *bridge.Element) []byte {

	// 发送消息到 rabbitmq
	fmt.Println(ele)
	// fmt.Println(body)
	// mybridge.Publish(ele)
	return nil
}

// 处理 客户端 关闭逻辑
func BeforeCloseHandler(clientid string) {
	manager.DeleteClient(clientid)
}

//Client 客户端（设备）
type Client struct {
	ID       string
	Conn     *SelfProto
	lifeSpan time.Duration // 生命周期
	Status   int           // 状态
	send     chan []byte   // 发送数据
}

// NewClient 创建客户端
func NewClient(key string, proto *SelfProto) *Client {
	return &Client{
		ID:       key,
		Conn:     proto,
		lifeSpan: 1 * time.Minute,
	}
}

// Loop 循环处理
func (c *Client) Loop() {
	for {
		select {
		case buf := <-c.send:
			logger.Debug("往设备写入数据", string(buf))
			c.Write(buf)

		case <-time.After(c.lifeSpan):
			logger.Debug("超时断开连接:", c.ID)
			c.Close()
		default:
			data, err := c.Conn.ReadMessage()
			if err != nil {
				//这里因为做了心跳，所以就没有加deadline时间，如果客户端断开连接
				//这里ReadByte方法返回一个io.EOF的错误，具体可考虑文档
				if err == io.EOF {
					logger.Error(fmt.Sprintf("client %s is close!\n", c.Conn.RemoteAddr().String()))
				}
				//在这里直接退出goroutine，关闭由defer操作完成
				logger.Error("Fail to decode error", err)
				c.Close()
			}
			go c.ProcessMessage(data)
		}

	}
}

//ProcessMessage 处理
func (c *Client) ProcessMessage(p []byte) {
	logger.Debugf("client %s ProcessMessage ", c.ID)
	// 读取内容
	ele, err := sdk.HandlePacket(p)
	if err != nil {
		logger.Errorf("client %s ProcessMessage error:%s", c.ID, err)
		return
	}

	// 基于消息内容 进行 相应的处理
	messageType := ele.MessageType
	switch messageType {
	case Connect:
		c.Register()
	case PubAck:
	case Publish:
		BusinessHandler(ele)
	case Heart:
		c.Send(Heart, 0, []byte(ele.Data))
	case Disconnect:
		c.Close()

	}
}

func (c *Client) Register() {}
func (c *Client) Ping()     { return }

//Close 关闭连接
func (c *Client) Close() {
	c.Conn.Close()
}

// Send send byte
func (c *Client) Send(messageType string, requestID int64, data []byte) error {
	timestamp := time.Now().Unix()
	if requestID == 0 {
		requestID = timestamp
	}
	var output = Package{
		Type:      messageType,
		RequestID: requestID,
		Timestamp: timestamp,
		Action:    "",
		Data:      data,
	}

	buf, _ := json.Marshal(output)

	return c.Conn.WriteMessage(buf)
}

//Write 写入连接数据
func (c *Client) Write(buf []byte) {
	//
	c.Conn.Write(buf)
}

// 连接管理相关的
//ConnManager 连接 管理器
type ConnManager struct {
	Clients sync.Map
}

// SetClient 存储
func (m *ConnManager) SetClient(key string, client *Client) {
	old, exist := m.Clients.Load(key)
	if exist {
		logger.Warn("client exist, close old...", key)
		ol, ok := old.(*Client)
		if ok {
			ol.Close()
		}
	}
	m.Clients.Store(key, client)
}

// GetClient 获取
func (m *ConnManager) GetClient(key string) *Client {
	value, ok := m.Clients.Load(key)
	if ok {
		return value.(*Client)
	}
	return nil
}

// DeleteClient 删除
func (m *ConnManager) DeleteClient(key string) {
	m.Clients.Delete(key)
}

// SelfProto 自定义协议包 SelfProtocol
// 原本 该功能 应该包含了 剥离基础协议的，但是设备上报的信息里，识别关键字在 接口路径上，故只能返回包含了 基础协议的数据内容
type SelfProto struct {
	conn net.Conn
}

func (p *SelfProto) ReadMessage() ([]byte, error) {
	buf, err := p.read()
	if err != nil {
		return nil, err
	}

	// 剥离 基础协议，重新处理后 再返回数据包
	return p.unpack(buf)
}

// 写入消息
func (p *SelfProto) WriteMessage(data []byte) error {
	// 封装 基础协议，封装完成后，再 发送数据包
	content := p.pack(data)

	// 写入 数据包
	_, err := p.Write(content)
	return err
}

func (p *SelfProto) Read(b []byte) (n int, err error) {
	return p.conn.Read(b)
}

func (p *SelfProto) Write(b []byte) (n int, err error) {
	return p.conn.Write(b)
}

func (p *SelfProto) Close() error {
	return p.conn.Close()
}

func (p *SelfProto) LocalAddr() net.Addr {
	return p.conn.LocalAddr()
}

func (p *SelfProto) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
}

func (p *SelfProto) SetDeadline(t time.Time) error {
	return p.conn.SetDeadline(t)
}

func (p *SelfProto) SetReadDeadline(t time.Time) error {
	return p.conn.SetReadDeadline(t)
}

func (p *SelfProto) SetWriteDeadline(t time.Time) error {
	return p.conn.SetWriteDeadline(t)
}

// 解包
func (p *SelfProto) unpack(data []byte) ([]byte, error) {
	return data, nil
}

// 封包
func (p *SelfProto) pack(data []byte) []byte {
	return data
}

func (p *SelfProto) read() ([]byte, error) {
	// 循环读取
	var err error
	myreader := bufio.NewReader(p.conn)
	// 建立一个只存4个的字节
	buf := make([]byte, 0)
	headendflag := "\r\n\r\n"
	for {
		// 按字节 读取
		b, err := myreader.ReadByte()
		if err != nil {
			return buf, err
		}
		buf = append(buf, b)
		last4Str := ""
		if len(buf) > 3 {
			last4bytes := buf[len(buf)-4:]
			last4Str = string(last4bytes)
		}
		flag := strings.EqualFold(headendflag, last4Str)
		if !flag {
			continue
		}
		// 匹配获取
		lengthStr := p.between(string(buf), "Content-Length: ", "\r\n")
		// fmt.Println(lengthStr)
		length, _ := strconv.Atoi(lengthStr)
		// 读取剩余的字段
		body := make([]byte, length)
		// n, err := myreader.Read(body)
		n, err := io.ReadFull(myreader, body)
		_ = n
		// fmt.Println(n)
		buf = append(buf, body...)
		// fmt.Println(string(buf))
		break
	}
	return buf, err
}

func (p *SelfProto) between(str string, start string, end string) string {
	// Get substring between two strings.
	s := strings.Index(str, start)
	if s == -1 {
		return ""
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return ""
	}
	return str[s : e+s]
}

func NewProto(conn net.Conn) *SelfProto {
	return &SelfProto{
		conn: conn,
	}
}

//InitTCP 启动tcp监听器
func InitTCP() {
	logger.Debug("tcp服务启动")
	var err error
	var l net.Listener
	tcpAddr = config.TCPAddr
	for {
		l, err = net.Listen("tcp", tcpAddr)
		logger.Info("Start Listening client on ", tcpAddr)

		if err != nil {
			logger.Error("Error listening on ", err)
			time.Sleep(1 * time.Second)
		} else {
			break // successfully listening
		}
	}
	tmpDelay := 10 * ACCEPT_MIN_SLEEP

	for {
		conn, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				logger.Error("Temporary Client Accept Error(%v), sleeping %dms", ne, tmpDelay/time.Millisecond)
				time.Sleep(tmpDelay)
				tmpDelay *= 2
				if tmpDelay > ACCEPT_MAX_SLEEP {
					tmpDelay = ACCEPT_MAX_SLEEP
				}
			} else {
				logger.Error("Accept error: %v", err)
			}
			continue
		}
		tmpDelay = ACCEPT_MIN_SLEEP
		go handleConn(conn)

	}
}

//handleConn 处理tcp连接
func handleConn(conn net.Conn) {
	logger.Info("处理新的 tcp 连接:", conn.RemoteAddr().String())
	defer conn.Close()

	protocol := NewProto(conn)

	// 首先读取 一个数据包
	buf, err := protocol.ReadMessage()
	if err != nil {
		//这里因为做了心跳，所以就没有加deadline时间，如果客户端断开连接
		//这里ReadByte方法返回一个io.EOF的错误，具体可考虑文档
		if err == io.EOF {
			logger.Error(fmt.Sprintf("client %s is close!\n", conn.RemoteAddr().String()))
		}
		//在这里直接退出goroutine，关闭由defer操作完成
		logger.Error("Fail to decode error", err)
		return
	}

	cid, err := sdk.ParseConnect(buf)
	logger.Info("建立连接的设备编号为:", cid)

	client := NewClient(cid, protocol)
	manager.SetClient(cid, client)

	client.Loop()
}

// HTTP 接口相关
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

	// 删除连接
	router.POST("api/faceguard", func(c *gin.Context) {
		// 解析 json 数据
		buf, _ := c.GetRawData()
		code, content, err := AsyncReq(buf)
		if err != nil {
			respondWithInfo(code, err.Error(), content, c)
		}
		respondWithInfo(0, "success", "", c)
	})

	httpAddr = config.HTTPAddr
	logger.Debug("http服务启动: ", httpAddr)
	router.Run(httpAddr)
}

//Sdk 处理设备功能模型等

// 定义消息结构
type Package struct {
	Type      string `json:"type"`       // 消息类型
	RequestID int64  `json:"request_id"` // 请求id
	Timestamp int64  `json:"timestamp"`  // 时间戳
	Action    string `json:"action"`     // 动作
	Data      []byte `json:"data"`       // 载荷
}

type Sdk struct{}

//_parse 解析包  topic payload
func (this *Sdk) _parse(buf []byte) (string, string, error) {
	eventType := "unknown"
	// 分割字符串
	str := string(buf)
	pos := strings.Index(str, "{")
	if pos < 1 {
		fmt.Println("未读取数据")
		return eventType, str, errors.New("未读取数据")
	}
	var input map[string]interface{}
	content := str[pos:]
	err := json.Unmarshal([]byte(content), &input)
	if err != nil {
		return eventType, content, errors.New("解析数据错误")
	}
	if strings.Contains(str, "HeartReportInfo") {
		// 处理心跳
		eventType = "heart"
	} else if strings.Contains(str, "PersonVerification") && !strings.Contains(str, "Notifications") {
		eventType = "open"
	} else if strings.Contains(str, "Response") {
		// 暂时不处理
		eventType = "response"
	}
	return eventType, content, nil
}

//ReadPacket 处理 消息, 返回 messageType, requestId , topic , payload
func (this *Sdk) ParsePacket(buf []byte) (string, int64, string, []byte, error) {
	topic := ""
	messageType := Unknown
	var requestid int64 = 0

	// 分割字符串
	str := string(buf)
	pos := strings.Index(str, "{")
	if pos < 1 {
		fmt.Println("未读取数据")
		return messageType, requestid, topic, nil, errors.New("未读取数据")
	}
	var input map[string]interface{}
	content := str[pos:]
	err := json.Unmarshal([]byte(content), &input)
	if err != nil {
		return messageType, requestid, topic, nil, errors.New("解析数据错误")
	}

	if strings.Contains(str, "HeartReportInfo") {
		return Heart, requestid, topic, nil, nil
	}

	if strings.Contains(str, "PersonVerification") && !strings.Contains(str, "Notifications") {
		messageType = Publish
		// 提取 requestid

	}

	if strings.Contains(str, "Response") {
		messageType = PubAck
		// 提取 requestid
	}
	return messageType, requestid, topic, nil, nil
}

//ParseConnect 解析连接信息
func (this *Sdk) ParseConnect(buf []byte) (string, error) {
	eventType, content, err := this._parse(buf)
	if err != nil {
		return "", err
	}
	if eventType != "heart" {
		return "", errors.New("不正确的数据包:" + eventType)
	}
	var input map[string]interface{}
	json.Unmarshal([]byte(content), &input)
	if deviceCode, ok := input["DeviceCode"]; ok {
		clientID := utils.String(deviceCode)
		return clientID, nil
	}
	return "", errors.New("设备编号不存在")
}

// Handle 传入消息包，返回: 消息类型， 消息id，topic, payload, error : 如果 是 publish 、 puback, 则进行其它处理
func (this *Sdk) HandlePacket(buf []byte) (*bridge.Element, error) {
	messageType, requestid, topic, body, err := this.ParsePacket(buf)
	if err != nil {
		return nil, err
	}

	ele := &bridge.Element{
		MessageType: messageType,
		RequestID:   requestid,
		Timestamp:   time.Now().Unix(),
	}

	if messageType == Heart {
		// 读取内容
		ele.Topic = topic
		ele.Data = this.RespPing()
		return ele, nil
	}

	if messageType == PubAck {
		// 处理 提取 关键 信息
		ele.Topic = topic
		ele.Data = string(body)
		return ele, nil

	}

	if messageType == Publish {
		ele.Topic = topic
		ele.Data = string(body)
		return ele, nil
	}

	return ele, errors.New("not match  messageType")
}

//RespPing 封装相应的响应
func (this *Sdk) RespPing() string {
	str := `{
	"ResponseURL": "/LAPI/V1.0/PACS/Controller/HeartReportInfo", 
	"Code": 0,
	"Data": {
		"Time": "%s" 
	}
}`
	str = fmt.Sprintf(str, utils.GetNormalTimeString(time.Now()))
	return str
}

func (this *Sdk) PackResp(body string, headers map[string]string) string {
	respHeader := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain;charset=ISO-8859-1\r\n"
	for key, val := range headers {
		respHeader = respHeader + key + ": " + val + "\r\n"
	}
	//拼装http返回的消息
	resp := respHeader + fmt.Sprintf("Content-Length: %d\r\n", len(body)) + "\r\n\r\n" + body
	return resp
}

func (this *Sdk) PackReq() {}

//Trans 处理 对应的 数据 提取
func (this *Sdk) Trans(topic string, data []byte) []byte {
	// 设备白名单
	if topic == "personAuthorized" {

	}

	// 删除白名单
	if topic == "personAuthorizedCancel" {

	}

	// 远程开门
	if topic == "deviceOpenFace" {

	}

	// 设备信息
	if topic == "deviceInfoFace" {

	}

	return nil
}

// 接口请求参数
type ReqParams struct {
	Topic   string
	Key     string
	Payload interface{}
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

	// 经过sdk转换为设备可以识别的请求
	buf, err = json.Marshal(input.Payload)
	fmt.Println("api接收到", string(buf))
	content := sdk.Trans(input.Topic, buf)
	cli.Write([]byte(content))

	// 生成消息 id
	msgId := utils.String(time.Now().Unix())
	logger.Debug("msgID :", msgId)

	// 建立消息通道
	ch := make(chan string, 1)
	msgChans[msgId] = ch

	// 循环读取 通道的响应，并增加超时退出
	var msg string
	select {
	case msg = <-ch:
		fmt.Println("data from channel #1", msg)
	// 如果把这个注释掉，则会阻塞 deadlock
	case <-time.After(5 * time.Second):
		fmt.Println("TimeOut")
		err = errors.New("响应超时TimeOut")
	}

	// 关闭通道，删除 key
	delete(msgChans, msgId)
	close(ch)

	// 返回结果
	return 0, msg, err
}
