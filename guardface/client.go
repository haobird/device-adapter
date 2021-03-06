package guardface

import (
	"context"
	"sync"
	"time"

	"github.com/haobird/logger"
)

const (
	Connected    = 1
	Disconnected = 2
)

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
		ID:         key,
		conn:       conn,
		keepalive:  lifeSpan,
		wmsgs:      make(chan []byte, 100),
		ctx:        ctx,
		cancelFunc: cancelFunc,
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
					// disconnectBuf, _ := json.Marshal(disconnectPack)
					// c.ProcessMessage(disconnectBuf)
					DisconnectHandler(c)
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
				// disconnectBuf, _ := json.Marshal(disconnectPack)
				// c.ProcessMessage(disconnectBuf)
				DisconnectHandler(c)
				return
			}
			// c.ProcessMessage(data)
			ProcessRawData(c, data)
		}

	}
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
