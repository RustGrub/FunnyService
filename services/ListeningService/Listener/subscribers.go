package Listener

import "github.com/RustGrub/FunnyGoService/consts"

func (s *ListeningService) SubscribeToNewGoods() error {
	err := s.msgBroker.Subscribe(consts.DefaultGoodsTopic, s.goodListener)
	if err != nil {
		return err
	}
	return nil
}
