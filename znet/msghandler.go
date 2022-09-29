package znet

import (
	"fmt"
	"strconv"

	"github.com/x627234313/zinx-test/utils"
	"github.com/x627234313/zinx-test/ziface"
)

type MsgHandle struct {
	// 存放每个msgID和对应的处理方法Handle
	Apis map[uint32]ziface.IRouter

	// 多任务工作池中worker数量
	WorkerPoolSize uint32
	// 多任务消息队列集合
	TaskQueue []chan ziface.IRequest
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (mh *MsgHandle) DoMsgHandle(r ziface.IRequest) {
	// 根据msgId获得对应的Handle
	handle, ok := mh.Apis[r.GetMsgId()]
	if !ok {
		fmt.Println("api msgId = ", r.GetMsgId(), " is NOT FOUND.")
		return
	}

	// 执行对应的方法
	handle.PreHandle(r)
	handle.Handle(r)
	handle.PostHandle(r)
}

func (mh *MsgHandle) AddRouter(msgid uint32, router ziface.IRouter) {
	// 首先判断msgid是否存在
	if _, ok := mh.Apis[msgid]; ok {
		panic("repeated api, msgId = " + strconv.Itoa(int(msgid)))
	}

	// 添加msgId对应的handle
	mh.Apis[msgid] = router
	fmt.Println("Add api msgId = ", msgid, " success.")
}

// 启动工作池
func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		workerChan := make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerChanReq)
		mh.TaskQueue[i] = workerChan
		go mh.startWorker(i)
	}

}

// 启动一个worker
func (mh *MsgHandle) startWorker(workeId int) {
	fmt.Println("Worker ID = ", workeId, " is started...")

	for {
		select {
		case req := <-mh.TaskQueue[workeId]:
			mh.DoMsgHandle(req)
		}
	}
}
