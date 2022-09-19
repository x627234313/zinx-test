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
}

type HandleFunc func(*net.TCPConn, []byte, int) error
