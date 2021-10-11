package guardface

import (
	"bufio"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/haobird/logger"
)

// ProtoConn 自定义协议包 ProtoConn
// 原本 该功能 应该包含了 剥离基础协议的，但是设备上报的信息里，识别关键字在 接口路径上，故只能返回包含了 基础协议的数据内容
type ProtoConn struct {
	conn net.Conn
}

func (p *ProtoConn) ReadMessage() ([]byte, error) {
	buf, err := p.read()
	if err != nil {
		return nil, err
	}

	// 剥离 基础协议，重新处理后 再返回数据包
	return p.unpack(buf)
}

// 写入消息
func (p *ProtoConn) WriteMessage(data []byte) error {
	// 封装 基础协议，封装完成后，再 发送数据包
	content := p.pack(data)

	// 写入 数据包
	logger.Debug("发送消息：", string(content))
	_, err := p.Write(content)
	return err
}

func (p *ProtoConn) Read(b []byte) (n int, err error) {
	return p.conn.Read(b)
}

func (p *ProtoConn) Write(b []byte) (n int, err error) {
	if p.conn != nil {
		return p.conn.Write(b)
	}
	return 0, errors.New("连接不存在")
}

func (p *ProtoConn) SetReadDeadline(t time.Time) error {
	return p.conn.SetReadDeadline(t)
}

func (p *ProtoConn) Close() error {
	return p.conn.Close()
}

// 解包
func (p *ProtoConn) unpack(data []byte) ([]byte, error) {
	return data, nil
}

// 封包
func (p *ProtoConn) pack(data []byte) []byte {
	return data
}

func (p *ProtoConn) read() ([]byte, error) {
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

func (p *ProtoConn) between(str string, start string, end string) string {
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

func NewProto(conn net.Conn) *ProtoConn {
	return &ProtoConn{
		conn: conn,
	}
}
