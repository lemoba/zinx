package main

import (
	"fmt"
	"github.com/lemoba/zinx/znet"
	"io"
	"net"
	"time"
)

/*
*
模拟客户端
*/
func main() {
	fmt.Println("client start...")

	time.Sleep(1 * time.Second)

	// 1. 直接连接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		// 发送封包的消息
		dp := znet.NewDataPack()

		binaryMsg, err := dp.Pack(znet.NewMsgPackge(0, []byte("ZinxV0.5 client Test Message")))
		if err != nil {
			fmt.Println("Pack message error: ", err)
			return
		}

		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("Send mssage error: ", err)
			return
		}

		// 从服务器得到回复
		// 1. 先读取流中的head部分 得到MsgID和dataLen
		binaryHead := make([]byte, dp.GetHeadLen())

		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("Read head error: ", err)
			break
		}

		// 将二进制的head拆包到msg结构体中
		msgHead, err := dp.Unpack(binaryHead)

		if err != nil {
			fmt.Println("Client unpack msgHead error: ", err)
		}

		// msg有数据
		if msgHead.GetMsgLen() > 0 {
			// 2. 在根据dataLen从data中读取数据
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msgHead.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("Read msg data error: ", err)
				return
			}

			fmt.Println("---> Receive Server Msg: ID = ", msg.Id,
				", Len = ", msg.DataLen,
				", Data = ", string(msg.Data))
		}

		// cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
