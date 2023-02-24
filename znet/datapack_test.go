package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPACK(t *testing.T) {
	/*
		模拟服务器
	*/
	// 1. 创建socketTCP
	listenner, err := net.Listen("tcp", "127.0.0.1:7777")

	if err != nil {
		fmt.Println("server listen err: ", err)
		return
	}

	// 创建一个go承载负责从客户端处理业务
	go func() {
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accept error: ", err)
			}

			go func(conn net.Conn) {
				// 处理客户端的请求
				// 拆包的过程
				// 定义一个拆包的对象
				dp := NewDataPack()

				for {
					// 1. 第一次从conn读， 把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error: ", err)
						break
					}

					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack error: ", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						// msg有数据，需要第二次读取
						// 2. 第二次从conn读， 根据head中的dataLen再读取data内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						// 根据datalen的长度再次从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data error: ", err)
							return
						}
						// 完整的一个消息已经读取完毕
						fmt.Println("---> Receive MsgID: ", msg.Id, ", datelen = ", msg.DataLen, ", data = ", string(msg.Data))
					}
				}
			}(conn)
		}
	}()
	/*
		模拟客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dail error: ", err)
		return
	}

	// 创建一个封包对象
	dp := NewDataPack()

	// 模拟粘包过程; 封装两个msg一同发送
	// 封装第一个msg包
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)

	if err != nil {
		fmt.Println("client pack msg1 error: ", err)
		return
	}
	// 封装第二个msg2包
	msg2 := &Message{
		Id:      2,
		DataLen: 5,
		Data:    []byte{'n', 'i', 'h', 'a', 'o'},
	}

	sendData2, err := dp.Pack(msg2)

	if err != nil {
		fmt.Println("client pack msg1 error: ", err)
		return
	}
	// 将两个包粘在一起
	sendData1 = append(sendData1, sendData2...)

	// 一次性发送给服务端
	conn.Write(sendData1)

	// 阻塞客户端
	select {}
}
