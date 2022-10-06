package znet

import (
	"errors"
	"fmt"
	"net"

	"github.com/x627234313/zinx-test/utils"
	"github.com/x627234313/zinx-test/ziface"
)

// 实例化一个IServer对象
type Server struct {
	IP        string
	Port      int
	IPVersion string
	Name      string

	// 该server的消息管理模块，把msgId和对应的业务处理方法绑定
	MsgHandle ziface.IMsgHandler

	// 当前 server 的连接管理器
	ConnMgr ziface.IConnMgr

	// 给当前 server 提供两个Hook方法的属性，参数就是conn 无返回值，在连接创建后、销毁前调用
	OnConnStart func(ziface.IConnection)
	OnConnStop  func(ziface.IConnection)
}

func CallBack(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("Conn Handle CallBack")

	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Printf("conn write error[%s]", err)
		return errors.New("CallBack Handle Error")
	}

	return nil
}

func (s *Server) Start() {
	fmt.Printf("[START] Zinx Server name=[%s], listenner ip=[%s], port=[%d] is starting.\n", s.Name, s.IP, s.Port)
	fmt.Printf("[ZINX] Version=[%s], MaxConn=[%d], MaxPacketSize=[%d], WorkerPooSize=[%d].\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize, utils.GlobalObject.WorkerPoolSize)

	go func() {
		tcpaddr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Printf("Resolve TCPAddr fail, err[%s]", err)
			return
		}

		tcplistener, err := net.ListenTCP(s.IPVersion, tcpaddr)
		if err != nil {
			fmt.Printf("Listen TCPAddr fail, err[%s], ip[%s], port[%d]", err, s.IP, s.Port)
			return
		}

		// 启动工作池
		s.MsgHandle.StartWorkerPool()

		var cid uint32

		for {
			tcpconn, err := tcplistener.AcceptTCP()
			if err != nil {
				fmt.Printf("Listener Accpet Fail. Error[%s]", err)
				continue
			}

			// 判断当前连接数是否超过最大连接数
			if s.ConnMgr.Count() >= utils.GlobalObject.MaxConn {
				fmt.Println("Too many connections")
				tcpconn.Close()
				continue
			}

			dealConn := NewConnection(s, tcpconn, cid, s.MsgHandle)
			cid++

			go dealConn.Start()

		}
	}()
}

func (s *Server) Stop() {
	// 当前 server 关闭时，清理所有的 connections
	fmt.Printf("[STOP] Zinx Server name=[%s] is stopped.", s.Name)

	s.ConnMgr.ClearAll()
}

func (s *Server) Serve() {
	s.Start()

	// 阻塞
	select {}
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandle.AddRouter(msgId, router)

	fmt.Println("Add Router success.")
}

func (s *Server) GetConnMgr() ziface.IConnMgr {
	return s.ConnMgr
}

// 实现注册、调用Hook函数的方法
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	s.OnConnStart(conn)
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	s.OnConnStop(conn)
}

func NewServer() ziface.IServer {
	utils.GlobalObject.Reload()

	return &Server{
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		IPVersion: "tcp4",
		Name:      utils.GlobalObject.Name,
		MsgHandle: NewMsgHandle(),
		ConnMgr:   NewConnMgr(),
	}
}
