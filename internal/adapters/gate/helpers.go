package gate

import (
	gate "github.com/gateio/gatews/go"

	"github.com/matrixbotio/exchange-gates-lib/internal/workers"
)

type gateSubsPayload struct {
	Channel string
	Payload []string
}

func getCandleSubsPayload(interval, pairSymbol string) gateSubsPayload {
	return gateSubsPayload{
		Channel: gateCandleChannel,
		Payload: []string{interval, pairSymbol},
	}
}

type gateUnsubscriber struct {
	srv  *gate.WsService
	data gateSubsPayload
}

func getUnsubscriber(srv *gate.WsService, data gateSubsPayload) workers.Unsubscriber {
	return &gateUnsubscriber{srv: srv, data: data}
}

func (u *gateUnsubscriber) Unsubscribe() error {
	return u.srv.UnSubscribe(u.data.Channel, u.data.Payload)
}
