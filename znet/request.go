package znet

import "zinx/ziface"

type Request struct {
	conn ziface.IConnection // 已经和客户端建立好的链接
	msg  ziface.IMessage    // 客户端请求的数据
}

// GetConnection 得到当前链接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetMsgID 得到请求的信息ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}

// GetMsgData 得到请求的信息数据
func (r *Request) GetMsgData() []byte {
	return r.msg.GetData()
}
