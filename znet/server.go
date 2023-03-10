package znet

import (
	"fmt"
	"github.com/lemoba/zinx/utils"
	"github.com/lemoba/zinx/ziface"
	"net"
)

// IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	// 服务器名称
	Name string
	// 服务器IP版本
	IPVersion string
	// 服务器监听IP
	IP string
	// 服务器监听端口
	Port int
	// 当前Server的消息管理模块，用来绑定MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandler
	// 该Server的连接管理器
	ConnManager ziface.IConnManager
	// 该Server创建连接之后自动调用Hook函数-OnConnStart
	OnConnStart func(connection ziface.IConnection)
	// 该Server销毁连接之前自动调用Hook函数-OnConnStop
	OnConnStop func(connection ziface.IConnection)
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s, listener at IP: %s, Port: %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version %s, MaxConn: %d, MaxPackageSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)

	go func() {
		// 0. 开启worker工作池以及消息队列
		s.MsgHandler.StartWorkerPool()

		// 1. 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("Resolve tcp addr error: ", err)
			return
		}
		// 2. 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)

		if err != nil {
			fmt.Println("Listen ", s.IPVersion, " err ", err)
			return
		}
		fmt.Println("Start Zinx server successful,", s.Name, "successful listening...")

		var cid uint32
		cid = 0

		// 3. 阻塞等待客户端连接，处理业务(读写)
		for {
			// 如果有客户端连接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error: ", err)
				continue
			}

			// 设置最大连接个数的判断，如果超过最大连接，那么则关闭此新的连接

			if s.ConnManager.Len() >= utils.GlobalObject.MaxConn {
				fmt.Println("Too Many Connections MaxConn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}
			// 处理新链接的业务方法和conn进行绑定，得到连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			// 启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server name", s.Name)
	s.ConnManager.ClearConn()
}

func (s *Server) Serve() {
	// 启动server的服务
	s.Start()

	// TODO 启动其他业务

	// 阻塞
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Successful!")
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnManager
}

/*
初始化Server模块的方法
*/
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandle(),
		ConnManager: NewConnManager(),
	}
	return s
}

// 注册onConnStart钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 注册onConnStop钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用onConnStart钩子函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

// 调用onConnStop钩子函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("Call OnConnStop()...")
		s.OnConnStop(conn)
	}
}
