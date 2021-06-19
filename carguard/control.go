package carguard

import "fmt"

// 控制接口
type Control interface {
	Publish(*Package) error
}

// mqtt 控制
func NewControl(name string) Control {
	switch name {
	case "mqtt":
		return InitMQTTControl()
	case "http":
		return InitMQTTControl()
	default:
		return &mockControl{}
	}
}

type mockControl struct{}

func (m *mockControl) Publish(packet *Package) error {
	topic := packet.Topic
	payload := packet.Data
	fmt.Println("topic:", topic)
	fmt.Println("payload:", string(payload))
	return nil
}
