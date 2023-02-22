package znet

import (
	"fmt"
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
	// 给当前服务器添加router
	Router ziface.IRouter
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listener at IP: %s, Port %d, is starting\n", s.IP, s.Port)

	go func() {
		// 1. 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}
		// 2. 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)

		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err ", err)
			return
		}
		fmt.Println("start Zinx server successful,", s.Name, "successful Listening...")

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

			// 处理新链接的业务方法和conn进行绑定，得到连接模块
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			// 启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// TODO implement me
	panic("implement me")
}

func (s *Server) Serve() {
	// 启动server的服务
	s.Start()

	// TODO 启动其他业务

	// 阻塞
	select {}
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router Successful!")
}

/*
初始化Server模块的方法
*/
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
		Router:    nil,
	}
	return s
}
