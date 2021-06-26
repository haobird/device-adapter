package guardface

// 数据包的几种状态
const (
	Unknown    = "unknown"
	Connect    = "connect"
	Heart      = "heart"
	Publish    = "publish"
	PubAck     = "puback"
	Command    = "command"
	Disconnect = "disconnect"
)

// 定义消息结构
type Package struct {
	MessageType string `json:"message_type"` // 消息类型
	RequestID   string `json:"request_id"`   // 请求id(时间戳)
	ClientID    string `json:"client_id"`    // 设备序列号
	Topic       string `json:"topic"`        // 主题
	Data        []byte `json:"data"`         // 载荷
}

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
