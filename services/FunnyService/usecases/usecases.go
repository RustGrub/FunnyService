package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/RustGrub/FunnyGoService/broker"
	"github.com/RustGrub/FunnyGoService/consts"
	"github.com/RustGrub/FunnyGoService/logger/std"
	"github.com/RustGrub/FunnyGoService/services/FunnyService/models"
	"github.com/RustGrub/FunnyGoService/services/FunnyService/repo/goods"
	"github.com/RustGrub/FunnyGoService/sql/txmanager"
	"github.com/RustGrub/FunnyGoService/subsidary"
	"sort"
	"time"
)

type FunnyService struct {
	txManager  txmanager.TxManager
	goodsRepo  goods.FsRepository
	goodsCache goods.Cache
	broker     broker.MessageBroker
}

func New(m txmanager.TxManager, r goods.FsRepository, c goods.Cache, b broker.MessageBroker) *FunnyService {
	return &FunnyService{
		txManager:  m,
		goodsRepo:  r,
		goodsCache: c,
		broker:     b,
	}
}

func (s *FunnyService) GetGood(ctx context.Context, goodID, projectID int) (res *models.Good, err error) {

	res = s.goodsCache.Get(ctx, goodID, projectID)
	if res != nil {
		return res, nil
	}

	err = s.txManager.Begin(ctx, func(ctx context.Context) error {
		res, err = s.goodsRepo.Get(ctx, goodID, projectID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if res != nil {
		rErr := s.goodsCache.Set(ctx, *res)
		if rErr != nil {
			subsidary.NewLogUsingContext(ctx, fmt.Sprintf("error set value in redis: %v", rErr), std.WarnignLevel)
		}
	}
	return res, nil
}

func (s *FunnyService) CreateGood(ctx context.Context, req models.CreateGoodRequest) (res *models.Good, err error) {
	err = s.txManager.Begin(ctx, func(ctx context.Context) error {
		res, err = s.goodsRepo.Create(ctx, req)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *FunnyService) UpdateGood(ctx context.Context, req models.UpdateGoodRequest) (res *models.Good, err error) {
	defer func() {
		// Условие ИЛИ для подстраховки, т.к. может быть паника
		if err != nil || res == nil {
			return
		}
		// Не обрабатываем ошибку, т.к. не влияет на дальнейшую работу
		if rErr := s.goodsCache.Set(ctx, *res); rErr != nil {
			subsidary.NewLogUsingContext(ctx, fmt.Sprintf("error set value in redis: %v", rErr), std.WarnignLevel)
		}

		if rErr := s.broker.PublishGoodAsLog(consts.DefaultGoodsTopic, models.GoodAsLog{
			GoodID:      res.GoodID,
			ProjectID:   res.ProjectID,
			Name:        res.Name,
			Description: res.Description,
			Priority:    res.Priority,
			Removed:     res.Removed,
			CreateDt:    time.Now(),
		}); rErr != nil {
			subsidary.NewLogUsingContext(ctx, fmt.Sprintf("new good not published in broker: %v", rErr), std.WarnignLevel)
		}
	}()

	err = s.txManager.Begin(ctx, func(ctx context.Context) error {
		res, err = s.goodsRepo.UpdateNameAndDescription(ctx, req)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *FunnyService) RemoveGood(ctx context.Context, req models.RemoveGoodRequest) (res models.RemoveGoodRequest, err error) {
	var good *models.Good
	err = s.txManager.Begin(ctx, func(ctx context.Context) error {
		good, err = s.goodsRepo.Remove(ctx, req)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return models.RemoveGoodRequest{}, err
	}

	// Можно делать во время транзакции выше, чтобы обеспечить согласованность в случае ошибки инвалидации кеша
	defer func() {
		// Условие ИЛИ для подстраховки, т.к. может быть паника
		if err != nil || good == nil {
			return
		}
		// Не обрабатываем ошибку, т.к. не влияет на дальнейшую работу
		if rErr := s.goodsCache.Set(ctx, *good); rErr != nil {
			subsidary.NewLogUsingContext(ctx, fmt.Sprintf("error set value in redis: %v", rErr), std.WarnignLevel)
		}

		if rErr := s.broker.PublishGoodAsLog(consts.DefaultGoodsTopic, models.GoodAsLog{
			GoodID:      good.GoodID,
			ProjectID:   good.ProjectID,
			Name:        good.Name,
			Description: good.Description,
			Priority:    good.Priority,
			Removed:     good.Removed,
			CreateDt:    time.Now(),
		}); rErr != nil {
			subsidary.NewLogUsingContext(ctx, fmt.Sprintf("new good not published in broker: %v", rErr), std.WarnignLevel)
		}
	}()
	return models.RemoveGoodRequest{
		Removed:   good.Removed,
		ProjectID: good.ProjectID,
		GoodID:    good.GoodID,
	}, nil
}

func (s *FunnyService) GetGoodsList(ctx context.Context, limit, offset int) (res models.GoodsListWithMeta, err error) {
	// Не стоит смотреть в кэше т.к. ttl - минута

	err = s.txManager.Begin(ctx, func(ctx context.Context) error {
		res.Goods, res.Meta.Removed, err = s.goodsRepo.GetListByLimitAndOffset(ctx, limit, offset)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return
	}
	res.Meta.Offset = offset
	res.Meta.Limit = limit
	res.Meta.Total = len(res.Goods)

	go func() {
		if res.Goods != nil {
			for _, v := range res.Goods {
				// Не сработает с ctx из запроса, т.к. тот схлопнется явно раньше, чем все запишется в кеш (унаследован от request.Context())
				rErr := s.goodsCache.Set(context.Background(), v)
				if rErr != nil {
					subsidary.NewLogUsingContext(context.Background(), fmt.Sprintf("error set value in redis: %v", rErr), std.WarnignLevel)
				}
			}
		}
	}()

	return
}

func (s *FunnyService) ReprioritizeGoods(ctx context.Context, req models.Reprioritize) (res models.ReprioritizeResponse, err error) {
	var good *models.Good
	var reprioritized []models.Good

	err = s.txManager.Begin(ctx, func(ctx context.Context) error {
		good, err = s.goodsRepo.Get(ctx, req.GoodID, req.ProjectID)
		if err != nil {
			return err
		}

		up := false
		if good.Priority < req.Priority {
			up = true
		} else if good.Priority == req.Priority {
			return errors.New("good is already at this priority")
		}

		reprioritized, err = s.goodsRepo.Reprioritize(ctx, req.Priority, good.Priority, up)
		if err != nil {
			return err
		}
		good, err = s.goodsRepo.UpdatePriority(ctx, good.GoodID, good.ProjectID, req.Priority)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return models.ReprioritizeResponse{}, err
	}

	reprioritized = append(reprioritized, *good)

	sort.Slice(res.Priorities, func(i, j int) bool {
		return res.Priorities[i].Priority < res.Priorities[j].Priority
	})
	res.Priorities = make([]models.ResetPriority, 0, len(reprioritized))
	// Формируем ответ клиенту, обновляем данные в кэше
	for _, v := range reprioritized {
		res.Priorities = append(res.Priorities, models.ResetPriority{
			GoodID:   v.GoodID,
			Priority: v.Priority,
		})

		if rErr := s.goodsCache.Set(ctx, v); rErr != nil {
			subsidary.NewLogUsingContext(ctx, fmt.Sprintf("error set value in redis: %v", rErr), std.WarnignLevel)
		}
		if rErr := s.broker.PublishGoodAsLog(consts.DefaultGoodsTopic, models.GoodAsLog{
			GoodID:      v.GoodID,
			ProjectID:   v.ProjectID,
			Name:        v.Name,
			Description: v.Description,
			Priority:    v.Priority,
			Removed:     v.Removed,
			CreateDt:    time.Now(),
		}); rErr != nil {
			subsidary.NewLogUsingContext(ctx, fmt.Sprintf("new good not published in broker: %v", rErr), std.WarnignLevel)
		}
	}

	return res, nil
}
