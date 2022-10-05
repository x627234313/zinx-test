package znet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/x627234313/zinx-test/ziface"
)

type ConnMgr struct {
	// 连接集合
	connections map[uint32]ziface.IConnection

	// 保护 map 的读写锁
	connLock sync.RWMutex
}

func NewConnMgr() *ConnMgr {
	return &ConnMgr{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (cm *ConnMgr) Add(conn ziface.IConnection) {
	// 保护共享资源，对 map 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 将 conn 添加到CommMgr中
	cm.connections[conn.GetConnId()] = conn

	fmt.Printf("ConnId = %d Add ConnMgr successfully. Current connection numbers = %d \n", conn.GetConnId(), cm.Count())
}

func (cm *ConnMgr) Remove(conn ziface.IConnection) error {
	// 保护共享资源，对 map 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 将 conn 从 ConnMgr 中移除
	if _, ok := cm.connections[conn.GetConnId()]; ok {
		delete(cm.connections, conn.GetConnId())
	} else {
		return errors.New("conn Not Found")
	}

	fmt.Printf("ConnId = %d Remove ConnMgr successfully. Current connection numbers = %d \n", conn.GetConnId(), cm.Count())
	return nil
}

func (cm *ConnMgr) GetConn(connId uint32) (ziface.IConnection, error) {
	// 保护共享资源，对 map 加读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	// 根据 connId 查找对应的 conn
	if c, ok := cm.connections[connId]; ok {
		return c, nil
	} else {
		return nil, errors.New("conn Not Found")
	}

}

func (cm *ConnMgr) Count() int {
	return len(cm.connections)
}

func (cm *ConnMgr) ClearAll() {
	// 保护共享资源，对 map 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 停止并清除所有的 conn
	for id, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, id)
	}

	fmt.Printf("Clear All Conn Successfully. Conn number = %d\n", cm.Count())
}
