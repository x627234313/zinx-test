package znet

import "github.com/x627234313/zinx-test/ziface"

type Request struct {
	connection ziface.IConnection
	msg        ziface.IMessage
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.connection
}

func (r *Request) GetData() []byte {
	return r.msg.GetMsgData()
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
