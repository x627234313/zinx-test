package znet

import (
	"errors"
	"fmt"
	"io"
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

	// 该connection的消息管理模块，把msgId和对应的业务处理方法绑定
	MsgHandle ziface.IMsgHandler
}

func NewConnection(conn *net.TCPConn, connId uint32, msgHandle ziface.IMsgHandler) *Connection {
	return &Connection{
		ConnId:    connId,
		Conn:      conn,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
		MsgHandle: msgHandle,
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

		//初始化一个拆包 封包对象
		dp := NewDataPack()

		// 根据headLen 读取head
		head := make([]byte, dp.GetHead())
		if _, err := io.ReadFull(c.GetTCPConn(), head); err != nil {
			fmt.Println("Read conn head error:", err)
			break
		}

		// 根据head 中datalen，读取data
		msg, err := dp.Unpack(head)
		if err != nil {
			fmt.Println("Unpack head error:", err)
			break
		}

		data := make([]byte, msg.GetMsgDataLen())
		if _, err := io.ReadFull(c.GetTCPConn(), data); err != nil {
			fmt.Println("Read conn data error:", err)
			break
		}

		msg.SetMsgData(data)

		fmt.Printf("Conn read msg id = %d, datalen = %d, data = %s\n", msg.GetMsgId(), msg.GetMsgDataLen(), string(msg.GetMsgData()))

		request := Request{
			connection: c,
			msg:        msg,
		}

		go c.MsgHandle.DoMsgHandle(&request)

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

func (c *Connection) GetTCPConn() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(id uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Conn is closed when sendmsg.")
	}

	// 初始化一个 封包对象
	dp := NewDataPack()

	// 对IMessage 进行封包
	binaryMsg, err := dp.Pack(NewMessage(id, data))
	if err != nil {
		fmt.Println("Pack message error, message id = ", id)
		return err
	}

	// 写回客户端
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write message id= ", id, "error")
		c.ExitChan <- true
		return err
	}

	return nil
}
