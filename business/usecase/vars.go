package usecase

import (
	"media-box-ui/adapter/broker"
)

type Broker interface {
	PublishState(data []byte)
	Subscribe(topic string, handler broker.MessageHandler)
}

type Image interface {
	Get(url string) (string, error)
}

var (
	stateTopic string
)

func SetStateTopic(t string) {
	stateTopic = t
}
