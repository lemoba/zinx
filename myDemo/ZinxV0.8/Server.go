package main

import (
	"fmt"
	"github.com/lemoba/zinx/ziface"
	"github.com/lemoba/zinx/znet"
)

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// test Handle
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Ping Router Handle")
	// 先读取客户端的数据，再回写ping

	fmt.Println("receive from client: msgId = ", request.GetMsgID(), " data = ", string(request.GetData()))

	if err := request.GetConnection().SendMsg(200, []byte("ping...ping...")); err != nil {
		fmt.Println(err)
	}
}

// hello test 自定义路由
type HelloRouter struct {
	znet.BaseRouter
}

// test Handle
func (hr *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Hello Router Handle")
	// 先读取客户端的数据，再回写ping

	fmt.Println("receive from client: msgId = ", request.GetMsgID(), " data = ", string(request.GetData()))

	if err := request.GetConnection().SendMsg(201, []byte("hello...hello...")); err != nil {
		fmt.Println(err)
	}
}

func main() {
	// 1. 创建一个server句柄， 使用Zinx的api
	s := znet.NewServer("[zinx V0.8]")

	// 2. 添加自定义router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	// 3. 启动server
	s.Serve()
}