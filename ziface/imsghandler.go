package ziface

// 定义消息管理模块，它有两个方法：
// 一个是执行msg对应的Handler
// 一个是添加msgId对应的Handler
type IMsgHandler interface {
	DoMsgHandle(IRequest)
	AddRouter(uint32, IRouter)

	// 启动工作池
	StartWorkerPool()
	// 将req发送给TaskQueue，由worker处理
	SendReqToTaskQueue(IRequest)
}
