package znet

import (
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

		for {
			tcpconn, err := tcplistener.AcceptTCP()
			if err != nil {
				fmt.Printf("Listener Accpet Fail. Error[%s]", err)
				continue
			}

			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := tcpconn.Read(buf)
					if err != nil {
						fmt.Println("TCPConn Read fail ", err)
						continue
					}

					if _, err := tcpconn.Write(buf[:cnt]); err != nil {
						fmt.Println("TCPConn Write fail ", err)
						continue
					}
				}

			}()

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
