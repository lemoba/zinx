package znet

import (
	"fmt"
	"github.com/lemoba/zinx/utils"
	"github.com/lemoba/zinx/ziface"
)

/*
消息处理模块的实现
*/

type MsgHandle struct {
	// 存在每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
	// 负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作Woker Pool的worker数量
	WorkerPoolSize uint32
}

// 初始化/创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
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

// 启动一个Worker工作池(开启一次)
func (mh *MsgHandle) StartWorkerPool() {
	// 根据wokerPoolSize分别开启Worker, 每个Worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 1. 当前的worker对应的channel消息队列开辟空间 第0个worker使用第0个channel ...
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 2. 启动当前的Worker，阻塞等待消息从channel传递进来
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (mh *MsgHandle) startOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started...")

	// 不断的阻塞等待对应消息队列的消息
	for {
		select {
		// 如果有消息过来
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1. 将消息平均分配给不同的worker
	// 根据客户端建立的ConnID来进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		"request MsgID = ", request.GetMsgID(), " to WorkerID = ", workerID)
	// 2. 将消息发送给对应的worker对应的TaskQueue即可
	mh.TaskQueue[workerID] <- request
}
