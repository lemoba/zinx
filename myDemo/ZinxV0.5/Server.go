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
	fmt.Println("Call Router Handle")
	// 先读取客户端的数据，再回写ping

	fmt.Println("receive from client: msgId = ", request.GetMsgID(), " data = ", string(request.GetData()))

	if err := request.GetConnection().SendMsg(0, []byte("ping...ping...")); err != nil {
		fmt.Println(err)
	}
}

func main() {
	// 1. 创建一个server句柄， 使用Zinx的api
	s := znet.NewServer("[zinx V0.5]")

	// 2. 添加自定义router
	s.AddRouter(&PingRouter{})

	// 3. 启动server
	s.Serve()
}
