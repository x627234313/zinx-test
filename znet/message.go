package znet

import "github.com/x627234313/zinx-test/ziface"

type Message struct {
	id      uint32
	data    []byte
	datalen uint32
}

func NewMessage(id uint32, data []byte) ziface.IMessage {
	return &Message{
		id:      id,
		data:    data,
		datalen: uint32(len(data)),
	}
}

func (m *Message) GetMsgId() uint32 {
	return m.id
}

func (m *Message) GetMsgData() []byte {
	return m.data
}

func (m *Message) GetMsgDataLen() uint32 {
	return m.datalen
}

func (m *Message) SetMsgId(id uint32) {
	m.id = id
}

func (m *Message) SetMsgData(data []byte) {
	m.data = data
}

func (m *Message) SetMsgDataLen(datalen uint32) {
	m.datalen = datalen
}
