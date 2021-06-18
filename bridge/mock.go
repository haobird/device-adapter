package bridge

import "fmt"

//mockMQ 模拟
type mockMQ struct{}

//Publish 发布
func (m *mockMQ) Publish(e *Element) error {
	fmt.Println(e)
	return nil
}
