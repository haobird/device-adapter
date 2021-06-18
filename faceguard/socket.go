package faceguard

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/haobird/logger"
)

const (
	ACCEPT_MIN_SLEEP = 100 * time.Millisecond
	ACCEPT_MAX_SLEEP = 10 * time.Second
)

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

	// 先执行注册的功能
	// msg := Message{
	// 	client: client,
	// 	packet: connectPack,
	// }
	// ProcessMessage(msg)
	// client.Register()
	client.Ping()

	client.Loop()
}
