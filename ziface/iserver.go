package ziface

type IServer interface {
	Start()
	Stop()
	Serve()

	// 增加AddRouter方法
	AddRouter(uint32, IRouter)

	// 获取连接管理器
	GetConnMgr() IConnMgr
}
