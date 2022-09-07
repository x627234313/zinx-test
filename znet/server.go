package znet

import (
	"errors"
	"fmt"
	"net"

	"github.com/x627234313/zinx-test/ziface"
)

// 实例化一个IServer对象
type Server struct {
	IP        string
	Port      int
	IPVersion string
	Name      string
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

		fmt.Printf("[Zin Server] IP:%s Port:%d Start Sucess, Listenning...\n", s.IP, s.Port)

		var cid uint32
		cid = 1

		for {
			tcpconn, err := tcplistener.AcceptTCP()
			if err != nil {
				fmt.Printf("Listener Accpet Fail. Error[%s]", err)
				continue
			}

			dealConn := NewConnection(tcpconn, cid, CallBack)
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

func NewServer(name string) ziface.IServer {
	return &Server{
		IP:        "0.0.0.0",
		Port:      9999,
		IPVersion: "tcp4",
		Name:      name,
	}
}
