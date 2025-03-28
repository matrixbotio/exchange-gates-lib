package workers

import "github.com/matrixbotio/exchange-gates-lib/pkg/structs"

type Unsubscriber interface {
	Unsubscribe() error
}

type channelsUnsubscriber struct {
	WsChannels structs.WorkerChannels
}

func CreateChannelsUnsubscriber(
	wsDone chan struct{},
	wsStop chan struct{},
) Unsubscriber {
	return &channelsUnsubscriber{
		WsChannels: structs.WorkerChannels{
			WsDone: wsDone,
			WsStop: wsStop,
		},
	}
}

func (s *channelsUnsubscriber) Unsubscribe() error {
	if s.WsChannels.WsStop == nil {
		return nil
	}

	go func() {
		if len(s.WsChannels.WsStop) == 0 {
			s.WsChannels.WsStop <- struct{}{}
		}
	}()
	return nil
}
