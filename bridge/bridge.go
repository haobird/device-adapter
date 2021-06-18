package bridge

//Element publish elements
type Element struct {
	MessageType string `json:"message_type"` // 消息类型
	RequestID   string `json:"request_id"`   // 请求id
	Timestamp   int64  `json:"timestamp"`    // 时间戳
	ClientID    string `json:"clientid"`     // 设备id
	Topic       string `json:"topic"`        // 路径或者主题或者表名
	Data        string `json:"data"`         // 详情
}

//Bridge 桥接接口
type Bridge interface {
	Publish(e *Element) error
}

//NewBridge 建立桥接
func NewBridge(name string) Bridge {
	switch name {
	case "rabbitmq":
		return InitBridgeMQ()
	default:
		return &mockMQ{}
	}
}
