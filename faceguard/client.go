package faceguard

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/haobird/logger"
)

const (
	Connected    = 1
	Disconnected = 2
)

var (
	connectPack = &Package{
		MessageType: Connect,
		RequestID:   "",
		ClientID:    "",
		Data:        nil,
	}

	disconnectPack = &Package{
		MessageType: Disconnect,
		RequestID:   "",
		ClientID:    "",
		Data:        nil,
	}
)

// 定义消息结构
type Package struct {
	MessageType string `json:"message_type"` // 消息类型
	RequestID   string `json:"request_id"`   // 请求id(时间戳)
	ClientID    string `json:"client_id"`    // 设备序列号
	Topic       string `json:"topic"`        // 主题
	Data        []byte `json:"data"`         // 载荷
}

//Client 客户端（设备）
type Client struct {
	ID                 string
	status             int // 状态
	mu                 sync.Mutex
	ctx                context.Context
	cancelFunc         context.CancelFunc
	conn               *ProtoConn
	keepalive          int         // 生命周期 多少秒
	wmsgs              chan []byte // 发送数据
	beforeCloseHandler func(string)
	msgHandler         func(*Package)
}

// NewClient 创建客户端
func NewClient(key string, conn *ProtoConn) *Client {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &Client{
		ID:                 key,
		conn:               conn,
		keepalive:          lifeSpan,
		wmsgs:              make(chan []byte, 100),
		beforeCloseHandler: BeforeCloseHandler,
		ctx:                ctx,
		cancelFunc:         cancelFunc,
	}
}

// 执行下发 任务
func (c *Client) dispatch() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case buf := <-c.wmsgs:
			logger.Debugf("[%s]往设备写入数据:%s", c.ID, string(buf))
			c.Write(buf)
		}
	}
}

// Loop 循环处理
func (c *Client) Loop() {
	nc := c.conn
	if nc == nil {
		return
	}

	// 执行 下发 任务
	go c.dispatch()

	keepAlive := time.Second * time.Duration(c.keepalive)
	timeOut := keepAlive + (keepAlive / 2)

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			//add read timeout 加deadline时间，如果客户端断开连接
			if keepAlive > 0 {
				if err := nc.SetReadDeadline(time.Now().Add(timeOut)); err != nil {
					logger.Errorf("[%s]set read timeout error: %s", c.ID, err)
					// msg := Message{
					// 	client: c,
					// 	packet: disconnectPack,
					// }
					// ProcessMessage(msg)
					disconnectBuf, _ := json.Marshal(disconnectPack)
					c.ProcessMessage(disconnectBuf)
					return
				}
			}
			data, err := nc.ReadMessage()
			if err != nil {
				logger.Errorf("[%s]read packet error: %s", c.ID, err)
				// msg := Message{
				// 	client: c,
				// 	packet: disconnectPack,
				// }
				disconnectBuf, _ := json.Marshal(disconnectPack)
				c.ProcessMessage(disconnectBuf)
				return
			}
			c.ProcessMessage(data)
			// SubmitWork(c, data)
		}

	}
}

func (c *Client) ProcessMessage(p []byte) {
	logger.Debugf("client %s ProcessMessage %s", c.ID, string(p))
	packet, err := sdk.HandlePacket(p)
	if err != nil {
		logger.Errorf("client %s ProcessMessage error:%s", c.ID, err)
		return
	}

	// 基于消息内容 进行 相应的处理
	messageType := packet.MessageType
	switch messageType {
	case Connect:
		logger.Infof("[%s] 设备Register ", c.ID)
		content := sdk.PackDeviceInfoReq()
		c.Write([]byte(content))
	case PubAck:
		BusinessHandler(packet)
	case Publish:
		BusinessHandler(packet)
	case Heart:
		c.Ping()
	case Disconnect:
		c.Close()
	}
}

func (c *Client) Register() {
	content := sdk.PackDeviceInfoReq()
	fmt.Println("Register 数据：", content)
	c.Write([]byte(content))
}

func (c *Client) Ping() {
	pong := sdk.Pong()
	c.Write([]byte(pong))
}

//Close 关闭连接
func (c *Client) Close() {
	if c.status == Disconnected {
		return
	}
	c.cancelFunc()

	c.status = Disconnected
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

//Write 写入连接数据
func (c *Client) Write(buf []byte) error {
	return c.conn.WriteMessage(buf)
}
