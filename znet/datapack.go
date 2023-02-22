package znet

import (
	"bytes"
	"encoding/binary"
	"github.com/lemoba/zinx/ziface"
)

// 封包，拆包的具体模块
type DataPack struct{}

// 拆包封包实例的一个初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	// Datalen uint32(4个字节) + ID uint32
	return 8
}

// 封包方法
// datalen|msgID|data
func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将dataLen写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// 将MsgId写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 将data数据写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包方法(将包的Head信息提取出来)
func (d *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// TODO implement me
	panic("implement me")
}
