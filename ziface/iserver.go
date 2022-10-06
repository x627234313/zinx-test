package ziface

type IServer interface {
	Start()
	Stop()
	Serve()

	// 增加AddRouter方法
	AddRouter(uint32, IRouter)

	// 获取连接管理器
	GetConnMgr() IConnMgr

	// 增加注册Hook函数的方法
	SetOnConnStart(func(IConnection))
	SetOnConnStop(func(IConnection))

	// 增加调用Hook函数的方法
	CallOnConnStart(IConnection)
	CallOnConnStop(IConnection)
}
