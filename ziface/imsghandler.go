package ziface

type IMsgHandler interface {
	DoMsgHandler(request IRequest)

	AddRouter(msgID uint32, router IRouter)

	// StartWorkerPool 启动worker工作池
	StartWorkerPool()

	// StartOneWorker 启动一个Worker工作流程
	StartOneWorker(int, chan IRequest)

	// SendMsgToTaskQueue 将消息交给TaskQueue,由worker进行处理
	SendMsgToTaskQueue(IRequest)
}
