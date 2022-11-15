package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
)

type MsgHandle struct {
	apis map[uint32]ziface.IRouter //存放每个MsgId 所对应的处理方法的map属性
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		apis: make(map[uint32]ziface.IRouter),
	}
}

// DoMsgHandler 马上以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgID(), " is not FOUND!")
		return
	}

	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.apis[msgID]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgID)))
	}
	//2 添加msg与api的绑定关系
	mh.apis[msgID] = router
	fmt.Println("Add api msgId = ", msgID)
}
