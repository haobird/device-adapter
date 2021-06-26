package guardcar

// 数据定义

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

// 定义消息包
type Package struct {
	MessageType string // 消息类型
	RequestID   string // 请求id(时间戳)
	ClientID    string // 设备序列号
	Action      string // 业务类型
	Topic       string // 主题
	Data        []byte // 载荷
}
