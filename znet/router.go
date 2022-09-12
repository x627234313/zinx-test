package znet

import "github.com/x627234313/zinx-test/ziface"

// 定义BaseRouter对象，实现IRouter接口，不实现具体方法，客户嵌入BaseRouter对象，重写方法
type BaseRouter struct{}

func (br *BaseRouter) PreHandle(ziface.IRequest) {}

func (br *BaseRouter) Handle(ziface.IRequest) {}

func (br *BaseRouter) PostHandle(ziface.IRequest) {}
