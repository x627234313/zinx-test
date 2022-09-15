package znet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/x627234313/zinx-test/utils"
	"github.com/x627234313/zinx-test/ziface"
)

type DataPack struct{}

func NewDataPack() ziface.IDataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHead() uint32 {
	// datalen uint32 (4字节) + id uint32 (4字节)
	return 8
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放byte的缓存
	dataBuff := bytes.NewBuffer([]byte{})

	// 写入 datalen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgDataLen()); err != nil {
		return nil, err
	}

	// 写入 id
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 写入 data
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个读取byte的io.Reader
	dataBuff := bytes.NewReader(binaryData)

	// 创建存放数据的Message，只解压 head，得到 datalen和id
	msg := &Message{}

	// 读取 datalen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.datalen); err != nil {
		return nil, err
	}

	// 读取 id
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.id); err != nil {
		return nil, err
	}

	// 判断消息中数据是否超过MaxPacketSize
	if msg.datalen > uint32(utils.GlobalObject.MaxPacketSize) && utils.GlobalObject.MaxPacketSize > 0 {
		return nil, errors.New("message data too large more than MaxPacketSize")
	}

	return msg, nil
}
