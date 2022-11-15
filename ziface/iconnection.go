package ziface

import "net"

type IConnection interface {
	// Start 启动
	Start()

	// Stop 停止链接 结束当前链接的工作
	Stop()

	// GetTcpConnect 获取当前链接的绑定socket conn
	GetTcpConnect() *net.TCPConn

	// GetConnID 获取当前链接模块的链接ID
	GetConnID() uint32

	// RemoteAddr 获取远程客户端的 TCP状态 IP port
	RemoteAddr() net.Addr

	// SendMsgData 发送数据 将数据发送给远程客户端
	SendMsgData(msgID uint32, data []byte) error
}

// HandleFunc 定义一个处理链接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
