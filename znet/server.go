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
	fmt.Printf("[START] Server name=[%s], listenner ip=[%s], port=[%d] is starting.\n", s.Name, s.IP, s.Port)
	fmt.Printf("[ZINX] Version=[%s], MaxConn=[%d], MaxPacketSize=[%d].\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)

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

		var cid uint32

		for {
			tcpconn, err := tcplistener.AcceptTCP()
			if err != nil {
				fmt.Printf("Listener Accpet Fail. Error[%s]", err)
				continue
			}

			dealConn := NewConnection(tcpconn, cid, s.MsgHandle)
			cid++

			go dealConn.Start()

		}
	}()
}

func (s *Server) Stop() {

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

func NewServer() ziface.IServer {
	utils.GlobalObject.Reload()

	return &Server{
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		IPVersion: "tcp4",
		Name:      utils.GlobalObject.Name,
		MsgHandle: NewMsgHandle(),
	}
}
