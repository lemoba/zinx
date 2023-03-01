package znet

import (
	"errors"
	"fmt"
	"github.com/lemoba/zinx/utils"
	"github.com/lemoba/zinx/ziface"
	"io"
	"net"
	"sync"
)

/*
 连接模块
*/

type Connection struct {
	// 当前Conn隶属于哪个Server
	TcpServer ziface.IServer
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 连接ID
	ConnID uint32
	// 当前连接的状态
	isClosed bool
	// 告知当前连接已经退出/停止的channel(由Reader告知Writer退出)
	ExitChan chan bool
	// 无缓冲的管道，用于读、写goroutine之间的消息通信
	msgChan chan []byte
	// 消息的管理MsgID和对应的处理业务API的关系
	MsgHandler ziface.IMsgHandler

	// 连接属性集合
	property map[string]interface{}
	// 保护连接属性读写锁
	propertyLock sync.RWMutex
}

// 初始化连接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		property:   make(map[string]interface{}),
	}

	// 将conn加入到ConnManager中
	c.TcpServer.GetConnManager().Add(c)
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running...]")
	defer fmt.Println("ConnID = ", c.ConnID, " [Reader is exit!], remote addr is ", c.RemoteAddr().String())
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
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制，将消息发送给worker工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 从路由中，找到注册绑定的Conn对应的router调用
			// 根据绑定好的MsgID 找到处理api对应的handler
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// 写消息的Goroutine, 专门发消息给客户端
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")
	defer fmt.Println(c.RemoteAddr().String(), " [conn Writer exit!]")

	// 不断地阻塞的等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error: ", err)
				return
			}
		case <-c.ExitChan:
			// 代表Reader已经退出，此时Writer也要退出
			return
		}
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
	c.msgChan <- binaryMsg

	return nil
}

func (c *Connection) Start() {
	fmt.Println("Connection Start... ConnID = ", c.ConnID)
	// 启动当前连接的读数据业务
	go c.StartReader()
	// 启动当前连接的写数据业务
	go c.StartWriter()

	// 按照开发者传递进来的，创建连接之后需要调用的处理业务，执行对应的Hook函数
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("Connection Stop... ConnID = ", c.ConnID)

	if c.isClosed == true {
		return
	}

	c.isClosed = true

	// 调用开发者注册的，销毁连接之前 需要执行的业务Hook函数
	c.TcpServer.CallOnConnStop(c)

	// 关闭socket连接
	c.Conn.Close()

	// 告知Writer关闭
	c.ExitChan <- true

	// 从ConnManager中移除当前连接
	c.TcpServer.GetConnManager().Remove(c)
	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, errors.New("not propery[" + key + "] found")
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
