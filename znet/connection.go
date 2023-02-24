package znet

import (
	"errors"
	"fmt"
	"github.com/lemoba/zinx/ziface"
	"io"
	"net"
)

/*
 连接模块
*/

type Connection struct {
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 连接ID
	ConnID uint32
	// 当前连接的状态
	isClosed bool
	// 告知当前连接已经退出/停止的channel
	ExitChan chan bool
	// 改连接处理的方法Router
	Router ziface.IRouter
}

// 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("ConnID = ", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中
		// buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		// _, err := c.Conn.Read(buf)
		// if err != nil {
		// 	fmt.Println("receive buf error ", err)
		// 	continue
		// }

		// 创建一个拆包解包对象
		dp := NewDataPack()

		// 读取客户端的Msg Head 8字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error: ", err)
			break
		}
		// 拆包,得到msgID和msgDatalen放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error: ", err)
		}
		// 根据datalen,再次读取Data,放入msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data err: ", err)
				break
			}
		}
		msg.SetData(data)
		// 得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		// 执行注册路由的方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
		// 从路由中，找到注册绑定的Conn对应的router调用
	}
}

// 将要发送给客户端的数据先封包再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}

	// 将data进行封包
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackge(msgId, data))
	if err != nil {
		fmt.Println("Pack message error: ", err, " MsgId = ", msgId)
		return errors.New("Pack error msg")
	}

	// 将数据发送给客户端
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write msg id ", msgId, " error: ", err)
		return errors.New("Conn Send error!")
	}
	return nil
}

func (c *Connection) Start() {
	fmt.Println("Connection Start... ConnID = ", c.ConnID)

	// 启动当前连接的读数据业务
	go c.StartReader()

	// TODO 启动当前连接的写数据业务
}

func (c *Connection) Stop() {
	fmt.Println("Connection Stop... ConnID = ", c.ConnID)

	if c.isClosed == true {
		return
	}

	c.isClosed = true

	// 关闭socket连接
	c.Conn.Close()

	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.GetConnID()
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
