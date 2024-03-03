package Listener

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RustGrub/FunnyGoService/services/FunnyService/models"
	_ "github.com/mailru/go-clickhouse/v2"
)

func (s *ListeningService) goodListener(msg []byte) {
	ctx := context.Background()
	fmt.Println("msg")

	var data models.GoodAsLog
	if err := json.Unmarshal(msg, &data); err != nil {
		s.logger.Warning(fmt.Sprintf("got message in broker but an error occured: %v", err))
		return
	}

	err := s.useCases.SaveGood(ctx, data)
	if err != nil {
		s.logger.Warning(fmt.Sprintf("an error occured while saving good %v in clickhouse: %v", data, err))
		return
	}
	return
}
