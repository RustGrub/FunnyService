package broker

import "github.com/RustGrub/FunnyGoService/services/FunnyService/models"

type SubscriptionFn = func(message []byte)

type MessageBroker interface {
	PublishGoodAsLog(topic string, good models.GoodAsLog) error
	Subscribe(topic string, subscription SubscriptionFn) error
}
