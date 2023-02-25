package znet

import "github.com/lemoba/zinx/ziface"

/*
消息处理模块的实现
*/

type MsgHandle struct {
	// 存在每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
}

// 初始化/创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

func (m *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// TODO implement me
	panic("implement me")
}

func (m *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// TODO implement me
	panic("implement me")
}
