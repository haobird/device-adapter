package faceguard

import (
	"deviceadapter/bridge"

	"github.com/haobird/goutils"
)

//RegisterHandler 注册
func RegisterHandler(c *Client) {
	content := sdk.PackDeviceInfoReq()
	c.Write([]byte(content))
}

// HeartHandler 心跳
func HeartHandler(c *Client) {
	pong := sdk.Pong()
	c.Write([]byte(pong))
}

func DisconnectHandler(c *Client) {
	manager.DeleteClient(c.ID)
	c.Close()
}

// PubAckHandler 设备响应的处理
func PubAckHandler(packet *Package) {
	requestID := packet.RequestID
	if ch, ok := msgChans[requestID]; ok {
		ch <- string(packet.Data)
	}
}

// PublishHandler 设备请求的处理
func PublishHandler(packet *Package) {
	// fmt.Println(body)
	ele := &bridge.Element{
		MessageType: packet.MessageType,
		RequestID:   packet.RequestID,
		Timestamp:   goutils.Int64(packet.RequestID),
		ClientID:    packet.ClientID,
		Topic:       packet.Topic,
		Data:        string(packet.Data),
	}
	mybridge.Publish(ele)
}
