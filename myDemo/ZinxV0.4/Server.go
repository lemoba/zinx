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

// test PreHandle
func (pr *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping error ", err)
	}
}

// test Handle
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...\n"))
	if err != nil {
		fmt.Println("call back ping... ping... error ", err)
	}
}

// test PostHandle
func (pr *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("call back after ping error ", err)
	}
}

func main() {
	// 1. 创建一个server句柄， 使用Zinx的api
	s := znet.NewServer("[zinx V0.4]")

	// 2. 添加自定义router
	s.AddRouter(&PingRouter{})

	// 3. 启动server
	s.Serve()
}
