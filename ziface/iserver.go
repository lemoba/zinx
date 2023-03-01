package ziface

// 定义服务器接口

type IServer interface {
	// 启动
	Start()
	// 停止
	Stop()
	// 运行服务器
	Serve()
	// 添加路由
	AddRouter(msgID uint32, router IRouter)
	// 获取server的连接管理
	GetConnManager() IConnManager
	// 注册onConnStart钩子函数的方法
	SetOnConnStart(func(connection IConnection))
	// 注册onConnStop钩子函数的方法
	SetOnConnStop(func(connection IConnection))
	// 调用onConnStart钩子函数的方法
	CallOnConnStart(connection IConnection)
	// 调用onConnStop钩子函数的方法
	CallOnConnStop(connection IConnection)
}
