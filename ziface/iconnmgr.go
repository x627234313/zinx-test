package ziface

// 定义连接管理模块
type IConnMgr interface {
	// 添加一个Conn
	Add(IConnection)

	// 移除一个Conn
	Remove(IConnection) error

	// 根据 id 查找 Conn
	GetConn(connId uint32) (IConnection, error)

	// 获取所有Conn的个数
	Count() int

	// 清理所有Conn
	ClearAll()
}
