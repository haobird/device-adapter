package bridge

import "fmt"

// 调用远程接口
type apiMQ struct {
	addr string
}

func (m *apiMQ) Publish(e *Element) error {
	fmt.Println(e)
	return nil
}
