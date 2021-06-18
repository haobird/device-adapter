package faceguard

import (
	"sync"

	"github.com/haobird/logger"
)

// 连接管理相关的
//ConnManager 连接 管理器
type ConnManager struct {
	Clients sync.Map
}

// SetClient 存储
func (m *ConnManager) SetClient(key string, client *Client) {
	old, exist := m.Clients.Load(key)
	if exist {
		logger.Warn("client exist, close old...", key)
		ol, ok := old.(*Client)
		if ok {
			ol.Close()
		}
	}
	m.Clients.Store(key, client)
}

// GetClient 获取
func (m *ConnManager) GetClient(key string) *Client {
	value, ok := m.Clients.Load(key)
	if ok {
		return value.(*Client)
	}
	return nil
}

// DeleteClient 删除
func (m *ConnManager) DeleteClient(key string) {
	m.Clients.Delete(key)
}
