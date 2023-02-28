package utils

import (
	"encoding/json"
	"github.com/lemoba/zinx/ziface"
	"os"
)

/*
存储一切有关Zinx框架的全局参数，供其他模块使用
一些参数可以通过zinx.json文件由用户进行配置
*/

type GlobalObj struct {
	/*
		Server
	*/
	TcpServer ziface.IServer // 当前Zinx全局的Server对象
	Host      string         // 当前服务器主机监听的IP
	TcpPort   int            // 当前服务器主机监听的端口号
	Name      string         // 当前服务器名称

	/*
		Zinx
	*/
	Version          string // 当前Zinx的版本号
	MaxConn          int    // 最大连接数
	MaxPackageSize   uint32 // 数据包最大值
	WorkerPoolSize   uint32 // 当前业务工作Worker池的Goroutine数量
	MaxWorkerTaskLen uint32 // Zinx框架允许用户最多开辟多少个Worker(限定条件)
}

/*
定义全局对外GlobalObj
*/
var GlobalObject *GlobalObj

/*
从zinx.json配置文件读取参数
*/
func (g *GlobalObj) LoadConfig() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// 将json文件数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
初始化当前GlobalObject
*/
func init() {
	// 如果配置文件没有加载,就使用下列默认值
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.5",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,   // worker工作池的队列个数
		MaxWorkerTaskLen: 1024, // 每个worker对应的消息队列的任务数量的最大值
	}

	// 从conf/zinx.json加载用户自定义参数
	GlobalObject.LoadConfig()
}
