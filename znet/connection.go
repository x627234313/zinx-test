package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/x627234313/zinx-test/utils"
	"github.com/x627234313/zinx-test/ziface"
)

type Connection struct {
	//  当前 conn 隶属于的 server
	TcpServer  ziface.IServer
	ConnId     uint32
	Conn       *net.TCPConn
	RemoteAddr net.Addr

	// 连接是否关闭
	isClosed bool

	// 连接关闭后向chan发送信息
	ExitChan chan bool

	// 无缓冲channel， 用于 读、写 goroutine 传递数据
	msgChan chan []byte

	// 该connection的消息管理模块，把msgId和对应的业务处理方法绑定
	MsgHandle ziface.IMsgHandler

	// 添加连接属性集合
	properity map[string]interface{}
	// 保护连接属性集合的锁
	properityLock sync.RWMutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connId uint32, msgHandle ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer: server,
		ConnId:    connId,
		Conn:      conn,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
		msgChan:   make(chan []byte),
		MsgHandle: msgHandle,
		properity: make(map[string]interface{}),
	}

	// 把当前 conn 添加到 ConnMgr 中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("[Conn Reader Goroutine is running]")
	// 下面语句会报错：panic: runtime error: invalid memory address or nil pointer dereference，可能是因为
	// c 实例化还未成功。
	//defer fmt.Printf("Conn[id=%d] Reader exit, remote addr is [%s]", c.ConnId, c.RemoteAddr.String())
	defer fmt.Println("[Conn Reader Goroutine exit.] ConnId = ", c.ConnId, " remote addr = ", c.GetRemoteAddr().String())
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

		if utils.GlobalObject.WorkerPoolSize > 0 {
			go c.MsgHandle.SendReqToTaskQueue(&request)
		} else {
			go c.MsgHandle.DoMsgHandle(&request)
		}

	}

}

func (c *Connection) StartWriter() {
	fmt.Println("[Conn Writer Goroutine is running]")
	defer fmt.Println("[Conn Writer Goroutine exit.] ConnId = ", c.ConnId, " remote addr = ", c.GetRemoteAddr().String())

	// 从msgChan中读取数据
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
		case <-c.ExitChan:
			return
		}
	}

}

func (c *Connection) Start() {
	// 启动从当前连接的读数据的业务
	go c.StartReader()

	// 启动从当前连接写数据的业务
	go c.StartWriter()

	fmt.Printf("Conn[id=%d] is start.\n", c.ConnId)

	// 连接创建完成之后，调用注册的Hook函数
	c.TcpServer.CallOnConnStart(c)

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

	// 连接关闭之前，调用注册的Hook函数
	c.TcpServer.CallOnConnStop(c)

	// 关闭socket连接
	c.Conn.Close()

	// 当前 conn 关闭时，从ConnMgr中移除
	c.TcpServer.GetConnMgr().Remove(c)

	// 告知Writer 关闭
	c.ExitChan <- true

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)

	fmt.Printf("Conn[id=%d] is stoped.\n", c.ConnId)
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
		return errors.New("Conn is closed when send msg.\n")
	}

	// 初始化一个 封包对象
	dp := NewDataPack()

	// 对IMessage 进行封包
	binaryMsg, err := dp.Pack(NewMessage(id, data))
	if err != nil {
		fmt.Println("Pack message error, message id = ", id)
		return err
	}

	//写回到msgChan中
	c.msgChan <- binaryMsg

	// 写回客户端
	/* if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write message id= ", id, "error")
		c.ExitChan <- true
		return err
	} */

	return nil
}

// 实现设置连接属性的方法
func (c *Connection) SetProperty(key string, value interface{}) {
	// 保护共享资源，加写锁
	c.properityLock.Lock()
	defer c.properityLock.Unlock()

	c.properity[key] = value
}

// 实现获取连接属性的方法
func (c *Connection) GetProperty(key string) (interface{}, error) {
	// 保护共享资源，加读锁
	c.properityLock.RLock()
	defer c.properityLock.RUnlock()

	if value, ok := c.properity[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("Conn Properity Not Found.")
	}
}

// 实现移除连接属性的方法
func (c *Connection) RemoveProperty(key string) {
	// 保护共享资源，加写锁
	c.properityLock.Lock()
	defer c.properityLock.Unlock()

	if _, ok := c.properity[key]; ok {
		delete(c.properity, key)
	}
}
