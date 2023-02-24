package ziface

import "net"

// 定义连接模块的抽象

type IConnection interface {
	// 启动连接， 当前连接开始工作
	Start()

	// 停止连接， 结束当前连接的工作
	Stop()

	// 获取当前连接的socket conn
	GetTCPConnection() *net.TCPConn

	// 获取当前连接的连接ID
	GetConnID() uint32

	// 获取远程客户端的TCP状态 IP Port
	RemoteAddr() net.Addr

	// 将解包后的数据发送给客户端
	SendMsg(msgId uint32, data []byte) error
}

// 定义一个处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
