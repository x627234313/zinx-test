package ziface

type IMessage interface {
	GetMsgId() uint32
	GetMsgData() []byte
	GetMsgDataLen() uint32

	SetMsgId(uint32)
	SetMsgData([]byte)
	SetMsgDataLen(uint32)
}
