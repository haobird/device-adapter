package guardcar

import "fmt"

type HTTPControl struct{}

func InitHTTPControl() *MQTTControl {
	return &MQTTControl{}
}

func (m *HTTPControl) Publish(packet *Package) error {
	topic := packet.Topic
	payload := packet.Data
	fmt.Println("topic:", topic)
	fmt.Println("payload:", string(payload))
	return nil
}
