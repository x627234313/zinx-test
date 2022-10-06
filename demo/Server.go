package main

import (
	"fmt"

	"github.com/x627234313/zinx-test/ziface"
	"github.com/x627234313/zinx-test/znet"
)

type TestRouter struct {
	znet.BaseRouter
}

func (tr *TestRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call TestRouter -> Handle")
	fmt.Println("recv from client : msgId=", request.GetMsgId(), "data=", string(request.GetData()))

	// 回写数据
	err := request.GetConnection().SendMsg(200, []byte("test...test..."))
	if err != nil {
		fmt.Println("Call TestRouter -> Handle Error")
		return
	}
}

type PingRouter struct {
	znet.BaseRouter
}

func (tr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter -> Handle")
	fmt.Println("recv from client : msgId=", request.GetMsgId(), "data=", string(request.GetData()))

	// 回写数据
	err := request.GetConnection().SendMsg(300, []byte("ping...ping..."))
	if err != nil {
		fmt.Println("Call PingRouter -> Handle Error")
		return
	}
}

func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("==> Call OnConnStart function: DoConnectionBegin.")

}

func DoConnectionClose(conn ziface.IConnection) {
	fmt.Println("==> Call OnConnStop function: DoConnectionClose.")
}

func main() {
	s := znet.NewServer()

	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionClose)

	s.AddRouter(0, &TestRouter{})
	s.AddRouter(1, &PingRouter{})

	s.Serve()
}
