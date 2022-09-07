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
	fmt.Println("Conn Reader Goroutine is running...")
	// 上面语句会报错：panic: runtime error: invalid memory address or nil pointer dereference，可能是因为
	// c 实例化还未成功。
	//defer fmt.Printf("Conn[id=%d] Reader exit, remote addr is [%s]", c.ConnId, c.RemoteAddr.String())
	defer fmt.Println(c.GetRemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	for {
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Printf("conn[id=%d] read error: %s", c.ConnId, err)
			c.ExitChan <- true
			continue
		}

		fmt.Printf("conn read: %s, cnt = %d\n", buf, cnt)

		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Printf("Conn[id=%d] handleAPI error[%s]", c.ConnId, err)
			c.ExitChan <- true
			return
		}
	}

}

func (c *Connection) Start() {
	go c.StartReader()

	fmt.Printf("Conn[id=%d] is start.\n", c.ConnId)

	for {
		select {
		case <-c.ExitChan:
			return
		}
	}
}

func (c *Connection) Stop() {
	if c.isClosed {
		fmt.Println("Conn already stoped.")
		return
	}

	c.isClosed = true
	c.Conn.Close()

	c.ExitChan <- true
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
