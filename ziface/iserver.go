package ziface

type IServer interface {
	Start()
	Stop()
	Serve()

	// 增加AddRouter方法
	AddRouter(uint32, IRouter)
}
