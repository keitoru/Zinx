package ziface

// IRequest 接口
type IRequest interface {
	// GetConnection 得到当前链接
	GetConnection() IConnection

	// GetMsgID 得到请求的信息ID
	GetMsgID() uint32

	// GetMsgData 得到请求的信息数据
	GetMsgData() []byte
}
