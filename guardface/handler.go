package guardface

import (
	"deviceadapter/bridge"

	"github.com/haobird/goutils"
)

//RegisterHandler 注册
func RegisterHandler(c *Client) {
	// 注册方法
	cid := c.ID
	manager.SetClient(cid, c)
	// 响应第一次的心跳
	HeartHandler(c)

	// 更新设备在线状态
	packet := sdk.Packet_deviceStatus(1)
	packet.ClientID = cid
	msg := Message{
		client: c,
		packet: packet,
	}
	ProcessMessage(msg)

	// 获取设备信息
	content := sdk.Packet_DeviceInfoReq()
	SubmitWork(c, content)
}

// HeartHandler 心跳
func HeartHandler(c *Client) {
	packet := sdk.Packet_heart()
	SubmitWork(c, packet)
}

// DisconnectHandler 断开
func DisconnectHandler(c *Client) {
	cid := c.ID
	manager.DeleteClient(cid)
	c.Close()

	// 更新设备在线状态
	packet := sdk.Packet_deviceStatus(0)
	packet.ClientID = cid
	msg := Message{
		client: c,
		packet: packet,
	}
	ProcessMessage(msg)
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
		MessageType: "entranceguard",
		RequestID:   packet.RequestID,
		Timestamp:   goutils.Int64(packet.RequestID),
		ClientID:    packet.ClientID,
		Topic:       packet.Topic,
		Data:        string(packet.Data),
	}
	mybridge.Publish(ele)
}
