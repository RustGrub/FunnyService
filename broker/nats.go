package broker

import (
	"github.com/RustGrub/FunnyGoService/config"
	"github.com/RustGrub/FunnyGoService/services/FunnyService/models"
	"github.com/nats-io/nats.go"
)

type NatsBroker struct {
	Nats *nats.Conn
}

func New(cfg *config.Config) (*NatsBroker, error) {
	nc, err := nats.Connect(cfg.Nats.Server + ":" + cfg.Nats.Port)
	if err != nil {
		return nil, err
	}
	return &NatsBroker{
		Nats: nc,
	}, nil
}

func (n *NatsBroker) PublishGoodAsLog(topic string, good models.GoodAsLog) error {
	data, err := good.MarshalBinary()
	if err != nil {
		return err
	}
	err = n.Nats.Publish(topic, data)
	if err != nil {
		return err
	}
	return nil
}
func (n *NatsBroker) Subscribe(topic string, subscription SubscriptionFn) error {
	_, err := n.Nats.Subscribe(topic, func(msg *nats.Msg) {
		subscription(msg.Data)
	})
	if err != nil {
		return err
	}
	return nil
}
