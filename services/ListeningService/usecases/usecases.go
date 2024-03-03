package usecases

import (
	"context"
	"github.com/RustGrub/FunnyGoService/logger"
	"github.com/RustGrub/FunnyGoService/services/FunnyService/models"
	"github.com/RustGrub/FunnyGoService/services/ListeningService/repo/goods"
)

type UseCases struct {
	goodsRepo goods.Repository
	logger.Logger
}

func New(l logger.Logger, r goods.Repository) *UseCases {
	return &UseCases{
		r,
		l,
	}
}

func (u *UseCases) SaveGood(ctx context.Context, good models.GoodAsLog) error {

	err := u.goodsRepo.Set(ctx, good)
	if err != nil {
		return err
	}
	return nil
}
