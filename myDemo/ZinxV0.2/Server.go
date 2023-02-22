package main

import "github.com/lemoba/zinx/znet"

func main() {
	// 1. 创建一个server句柄， 使用Zinx的api
	s := znet.NewServer("[zinx V0.2]")
	// 启动server
	s.Serve()
}
