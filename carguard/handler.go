package carguard

//RegisterHandler 注册
func RegisterHandler(clientID string) {

}

// HeartHandler 心跳
func HeartHandler(clientID string) {

}

// DisconnectHandler 断开
func DisconnectHandler(clientID string) {
	cache.Delete(clientID)
}
