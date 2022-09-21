package znet

import (
	"fmt"
	"strconv"

	"github.com/x627234313/zinx-test/ziface"
)

type MsgHandle struct {
	// 存放每个msgID和对应的处理方法Handle
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
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
