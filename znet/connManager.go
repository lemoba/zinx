package znet

import (
	"errors"
	"fmt"
	"github.com/lemoba/zinx/ziface"
	"sync"
)

/*
连接管理模块
*/

type ConnManager struct {
	connections map[uint32]ziface.IConnection // 管理的连接集合
	connLock    sync.RWMutex                  // 读写锁
}

// 创建当前连接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (cg *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源map, 加写锁
	cg.connLock.Lock()
	defer cg.connLock.Unlock()

	// 将conn加入到ConnManager中
	cg.connections[conn.GetConnID()] = conn
	fmt.Println("ConnID = ", conn.GetConnID(), " Add to ConnManager successful: conn num = ", cg.Len())
}

func (cg *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源map, 加写锁
	cg.connLock.Lock()
	defer cg.connLock.Unlock()

	// 删除连接信息
	delete(cg.connections, conn.GetConnID())
	fmt.Println("ConnID = ", conn.GetConnID(), "Remove from ConnManager successful: conn num = ", cg.Len())
}

func (cg *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源map, 加读锁
	cg.connLock.RLock()
	defer cg.connLock.RUnlock()

	if conn, ok := cg.connections[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("Connection not found!")
}

func (cg *ConnManager) Len() int {
	return len(cg.connections)
}

func (cg *ConnManager) ClearConn() {
	cg.connLock.Lock()
	defer cg.connLock.Unlock()

	// 删除conn并停止conn的工作
	for connID, conn := range cg.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(cg.connections, connID)
	}

	fmt.Println("Clear All connections successful! conn nums = ", cg.Len())
}
