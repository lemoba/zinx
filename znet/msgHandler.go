package znet

import (
	"fmt"
	"github.com/lemoba/zinx/ziface"
)

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

func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 1. 从request找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("Api msgID = ", request.GetMsgID(), " is not found! need register!!")
		return
	}
	// 2. 根据MsID调度对应的router
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1. 判断当前msg绑定的API处理方法是否存在
	if _, ok := mh.Apis[msgID]; ok {
		fmt.Println("Repeat api, msgID = ", msgID)
		return
	}
	// 2. 添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, " successful!")
}
