package ziface

// 封包 拆包 模块
// 解决TCP 数据传输中粘包的用户
// 使用TLV(Type-Len-Value)打包格式：datalen + msgid + data
type IDataPack interface {
	// 获取消息头
	GetHead() uint32

	// 把IMessage 打包成字节流
	Pack(IMessage) ([]byte, error)

	// 把字节流解析成IMessage
	Unpack([]byte) (IMessage, error)
}
