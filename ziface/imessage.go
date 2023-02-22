package ziface

/*
将请求消息封装到一个message中
*/
type IMessage interface {
	GetMsgId() uint32   // 获取消息ID
	GetDataLen() uint32 // 获取消息长度
	GetData() []byte    // 获取消息内容

	SetMsgId(uint32)   // 设置消息ID
	SetDataLen(uint32) // 设置消息长度
	SetData([]byte)    // 设置消息内容
}
