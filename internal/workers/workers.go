package workers

import structs "github.com/matrixbotio/exchange-gates-lib/pkg/structs"

type workerBase struct {
	WsChannels *structs.WorkerChannels
}

func (w *workerBase) Stop() {
	if w.WsChannels == nil || w.WsChannels.WsStop == nil {
		return
	}

	go func() {
		if len(w.WsChannels.WsStop) == 0 {
			w.WsChannels.WsStop <- struct{}{}
		}
	}()
}
