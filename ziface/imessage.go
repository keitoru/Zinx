package ziface

/*
	将请求的一个消息封装到message中，定义抽象层接口
*/
type IMessage interface {
	// GetMsgId 获取消息ID
	GetMsgId() uint32

	// GetMsgLen 获取消息数据段长度
	GetMsgLen() uint32

	// GetData 获取消息内容
	GetData() []byte

	// SetMsgId c设置消息ID
	SetMsgId(uint32)

	// SetMsgLen 设置消息数据段长度
	SetMsgLen(uint32)

	// SetData 设置消息内容
	SetData([]byte)
}
