package ziface

import "net"

type IConnection interface {
	Start()
	Stop()

	// 获得连接ID
	GetConnId() uint32

	// 获得sock套接字
	GetTCPConn() *net.TCPConn

	// 获得客户端IP、端口
	GetRemoteAddr() net.Addr

	// 发送数据
	SendMsg(id uint32, data []byte) error

	// 设置连接属性
	SetProperty(string, interface{})
	// 获取连接属性
	GetProperty(string) (interface{}, error)
	// 移除连接属性
	RemoveProperty(string)
}

type HandleFunc func(*net.TCPConn, []byte, int) error
