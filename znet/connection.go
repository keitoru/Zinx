package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

// Connection 链接模块
type Connection struct {
	Conn       *net.TCPConn       // 当前链接的socket TCP套接字
	ConnId     uint32             // 链接ID
	IsClose    bool               // 当前的链接状态
	MsgHandler ziface.IMsgHandler //消息管理MsgId和对应处理方法的消息管理模块
	ExitChan   chan bool          // 告知当前业务已经退出/停止的channel
	MsgChan    chan []byte
}

// NewConnection 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnId:     connID,
		IsClose:    false,
		MsgHandler: msgHandler,
		ExitChan:   make(chan bool, 1),
		MsgChan:    make(chan []byte),
	}

	return c
}

// StartReader 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")

	defer fmt.Println("[Reader is exit], connID=", c.ConnId, "remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中， 最大512字节
		//buf := make([]byte, utils.GlobalObject.MaxPacketSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err", err)
		//	continue
		//}

		// 创建拆包解包的对象
		dp := NewDataPack()

		// 读取客户端的Msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTcpConnect(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			break
		}

		//拆包，得到msgid 和 datalen 放在msg中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			break
		}

		//根据 dataLen 读取 data，放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTcpConnect(), data); err != nil {
				fmt.Println("read msg data error ", err)
				break
			}
		}
		msg.SetData(data)

		//得到当前客户端请求的Request数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经启动工作池机制，将消息交给Worker处理
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//从绑定好的消息和对应的处理方法中执行对应的Handle方法
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

// StartWriter 写消息Goroutine，专门发送消息给客户端的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")

	defer fmt.Println("[conn Writer exit!]", c.RemoteAddr().String())

	for {
		select {
		case data := <-c.MsgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err)
				return
			}
		case <-c.ExitChan:
			//conn已经关闭
			return
		}
	}
}

// Start 启动链接 让当前链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start... ConnID=", c.ConnId)
	// 启动从当前链接的读数据业务
	go c.StartReader()

	// 启动从当前链接的写数据业务
	go c.StartWriter()
}

// Stop 停止链接 结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop... ConnID=", c.ConnId)

	// 如果当前链接已关闭
	if c.IsClose == true {
		return
	}
	c.IsClose = true

	// 关闭socket链接
	if err := c.Conn.Close(); err != nil {
		fmt.Println("conn close err:", err)
	}

	// 告知Write关闭
	c.ExitChan <- true

	// 回收资源
	close(c.ExitChan)
	close(c.MsgChan)
}

// GetTcpConnect 获取当前链接的绑定socket conn
func (c *Connection) GetTcpConnect() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnId
}

// RemoteAddr 获取远程客户端的 TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsgData 发送数据 将数据发送给远程客户端
func (c *Connection) SendMsgData(msgID uint32, data []byte) error {
	if c.IsClose == true {
		return errors.New("Connection closed when send msg")
	}

	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgID)
		return errors.New("Pack error msg ")
	}

	// 将数据发送给客户端
	//if _, err = c.Conn.Write(binaryMsg); err != nil {
	//	fmt.Println("Write msg id ", msgID, " error ")
	//	return errors.New("conn Write error")
	//}
	c.MsgChan <- binaryMsg

	return nil
}
