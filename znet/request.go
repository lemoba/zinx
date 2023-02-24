package znet

import "github.com/lemoba/zinx/ziface"

type Request struct {
	// 已经和客户端建立连接的map
	conn ziface.IConnection
	// 客户端请求的数据
	msg ziface.IMessage
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// 获取请求的消息数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// 获取请求消息的ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
