package znet

import (
	"fmt"
	"net"

	"github.com/x627234313/zinx-test/ziface"
)

type Connection struct {
	ConnId     uint32
	Conn       *net.TCPConn
	RemoteAddr net.Addr

	// 连接是否关闭
	isClosed bool

	// 连接关闭后向chan发送信息
	ExitChan chan bool

	// 连接处理业务的方法
	handleAPI ziface.HandleFunc
}

func NewConnection(conn *net.TCPConn, connId uint32, callback ziface.HandleFunc) *Connection {
	return &Connection{
		ConnId:    connId,
		Conn:      conn,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
		handleAPI: callback,
	}
}

func (c *Connection) StartReader() {

}

func (c *Connection) Start() {
	fmt.Printf("Reader Goroutine [id=%d] is running", c.ConnId)

	go c.StartReader()
}

func (c *Connection) Stop() {
	if c.isClosed {
		fmt.Println("Conn already closed.")
		return
	}

	c.isClosed = true
	c.Conn.Close()

	close(c.ExitChan)

	fmt.Printf("Conn[id=%d] is stoped", c.ConnId)
}

func (c *Connection) GetConn() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	return nil
}
