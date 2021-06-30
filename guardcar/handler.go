package guardcar

import (
	"deviceadapter/bridge"

	"github.com/haobird/goutils"
)

//RegisterHandler 注册
func RegisterHandler(clientID string) {

}

// HeartHandler 心跳
func HeartHandler(clientID string) {

}

// DisconnectHandler 断开
func DisconnectHandler(clientID string) {
	cache.Delete(clientID)

	// 更新设备在线状态
	packet := things.Packet_deviceStatus(0)
	packet.ClientID = clientID

	// 执行上报逻辑
	ele := &bridge.Element{
		MessageType: "vehicledetection",
		RequestID:   packet.RequestID,
		Timestamp:   goutils.GetTime().Unix(),
		ClientID:    packet.ClientID,
		Topic:       packet.Topic,
		Data:        string(packet.Data),
	}
	mybridge.Publish(ele)
}
