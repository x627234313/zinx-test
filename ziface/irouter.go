package ziface

// 定义Router 抽象对象，处理对象是IRequest
// 有业务处理之前的Hook方法
// 业务处理的主方法
// 业务处理之后的Hook方法
type IRouter interface {
	PreHandle(IRequest)
	Handle(IRequest)
	PostHandle(IRequest)
}
