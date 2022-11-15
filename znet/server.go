package znet

import (
	"fmt"
	"net"
	"time"
	"zinx/utils"
	"zinx/ziface"
)

// iServer 接口实现，定义一个Server服务类
type Server struct {
	//服务器的名称
	Name string
	//tcp4 or other
	IPVersion string
	//服务绑定的IP地址
	IP string
	//服务绑定的端口
	Port int
	//当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
	MsgHandler ziface.IMsgHandler
	//无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgChan chan []byte
}

//============== 实现 ziface.IServer 里的全部接口方法 ========

// Start 开启网络服务
func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)

	go func() {
		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("Resolve TCP Addr error:", err)
			return
		}

		//2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen tcp error:", err)
			return
		}

		//已经监听成功
		fmt.Println("start Zinx server  ", s.Name, " success, now listening...")
		var cid uint32

		//3 启动server网络连接业务
		for {
			//3.1 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}

			// 将处理新链接的业务方法 和 conn 进行绑定，得到我们的链接模块
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++

			// 启动当前的业务处理
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)
}

func (s *Server) Serve() {
	s.Start()

	//阻塞,否则主Go退出， listener的go将会退出
	for {
		time.Sleep(10 * time.Second)
	}
}

// AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router success!")
}

// NewServer 创建一个服务器句柄
func NewServer(name string) ziface.IServer {
	//先初始化全局配置文件
	utils.GlobalObject.Reload()

	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(), //msgHandler 初始化
		msgChan:    make(chan []byte),
	}

	return s
}
