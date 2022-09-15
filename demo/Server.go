package main

import (
	"fmt"

	"github.com/x627234313/zinx-test/ziface"
	"github.com/x627234313/zinx-test/znet"
)

type TestRouter struct {
	znet.BaseRouter
}

func (tr *TestRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call TestRouter -> PreHandle :")
	_, err := request.GetConnection().GetTCPConn().Write([]byte("before test...\n"))
	if err != nil {
		fmt.Println("Call TestRouter -> PreHandle Error")
		return
	}
}

func (tr *TestRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call TestRouter -> Handle :")
	_, err := request.GetConnection().GetTCPConn().Write([]byte("test...\n"))
	if err != nil {
		fmt.Println("Call TestRouter -> Handle Error")
		return
	}
}

func (tr *TestRouter) PostHandle(request ziface.IRequest) {
	fmt.Println(" Call TestRouter -> PostHandle :")
	_, err := request.GetConnection().GetTCPConn().Write([]byte("after test...\n"))
	if err != nil {
		fmt.Println("Call TestRouter -> PostHandle Error")
		return
	}
}

func main() {
	s := znet.NewServer()

	s.AddRouter(&TestRouter{})

	s.Serve()
}
