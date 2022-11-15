package znet

type Message struct {
	MsgId  uint32 //消息的ID
	MsgLen uint32 //消息的长度
	Data   []byte //消息的内容
}

func NewMsgPackage(msgId uint32, data []byte) *Message {
	m := &Message{
		MsgId:  msgId,
		MsgLen: uint32(len(data)),
		Data:   data,
	}

	return m
}

// GetMsgId 获取消息ID
func (m *Message) GetMsgId() uint32 {
	return m.MsgId
}

// GetMsgLen 获取消息数据段长度
func (m *Message) GetMsgLen() uint32 {
	return m.MsgLen
}

// GetData 获取消息内容
func (m *Message) GetData() []byte {
	return m.Data
}

// SetMsgId c设置消息ID
func (m *Message) SetMsgId(msgId uint32) {
	m.MsgId = msgId
}

// SetMsgLen 设置消息数据段长度
func (m *Message) SetMsgLen(len uint32) {
	m.MsgLen = len
}

// SetData 设置消息内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}
