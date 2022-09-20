package ziface

type IMsgHandler interface {
	DoMsgHandle(IRequest)
	AddRouter(uint32, IRouter)
}
