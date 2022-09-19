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
	err := request.GetConnection().SendMsg(1, []byte("test...test..."))
	if err != nil {
		fmt.Println("Call TestRouter -> Handle Error")
		return
	}
}

func main() {
	s := znet.NewServer()

	s.AddRouter(&TestRouter{})

	s.Serve()
}
