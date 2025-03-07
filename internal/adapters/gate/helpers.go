package gate

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
