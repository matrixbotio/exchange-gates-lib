package workers

// WorkerChannels - channels container to control the worker
type WorkerChannels struct {
	WsDone chan struct{}
	WsStop chan struct{}
}
