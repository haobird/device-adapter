package bridge

import (
	"encoding/json"

	"github.com/haobird/goutils"
)

// 调用远程接口
type apiMQ struct {
	addr string
}

func InitApiMQ() *apiMQ {
	return &apiMQ{
		addr: "http://iot-community.b.mi.com/api/v1/",
	}
}

func (m *apiMQ) Publish(e *Element) error {
	kind := e.MessageType
	url := m.addr + kind
	topic := e.Topic
	key := e.ClientID
	data := e.Data
	var payload map[string]interface{}
	json.Unmarshal([]byte(data), &payload)
	var req = map[string]interface{}{
		"key":     key,
		"topic":   topic,
		"payload": payload,
	}

	buf, _ := json.Marshal(req)
	header := map[string]string{
		"Content-Type": "application/json;charset=UTF-8",
	}
	goutils.Request(url, "POST", buf, header)

	return nil
}
