package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	// 1 创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("Server listen err:", err)
		return
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Server accept error:", err)
		}

		go func(conn net.Conn) {
			// 处理客户端的请求
			// ----> 拆包 <----
			// 定义一个拆包对象
			dp := NewDataPack()
			for {
				// 1 第一次从conn读， 把包的head读出来
				headData := make([]byte, dp.GetHeadLen())
				if _, err := io.ReadFull(conn, headData); err != nil {
					fmt.Println("read head error")
					return
				}

				msgHead, err := dp.UnPack(headData)
				if err != nil {
					fmt.Println("Server unpack err:", err)
					return
				}

				if msgHead.GetMsgLen() > 0 {
					// 2 第二次从conn读，根据head中的dataLen，再读取data的内容
					msg := msgHead.(*Message)
					msg.Data = make([]byte, msg.GetMsgLen())

					// 根据msgLen的长度再次从io流中读取
					if _, err := io.ReadFull(conn, msg.Data); err != nil {
						fmt.Println("Server unpack data err:", err)
						return
					}

					fmt.Println("---> Recv MsgID:", msg.MsgId, ", MsgLen = ", msg.MsgLen, ", data = ", string(msg.Data))
				}

			}
		}(conn)
	}()

	// 模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}

	// 创建一个封包对象 dp
	dp := NewDataPack()

	// 模拟粘包过程，封装两个msg一同发送
	// 封装第一个msg1包
	msg1 := &Message{
		MsgId:  1,
		MsgLen: 4,
		Data:   []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("Client pack msg1 err:", err)
		return
	}

	// 封装第二个msg2包
	msg2 := &Message{
		MsgId:  2,
		MsgLen: 7,
		Data:   []byte{'h', 'e', 'l', 'l', 'o', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("Client pack msg2 err:", err)
		return
	}

	// 将两个包黏在一起
	sendData1 = append(sendData1, sendData2...)

	// 一次性发给服务器
	if _, err := conn.Write(sendData1); err != nil {
		fmt.Println("Client conn write err:", err)
		return
	}

	// 客户端阻塞
	select {}

}
