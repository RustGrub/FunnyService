package goods

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RustGrub/FunnyGoService/cache"
	"github.com/RustGrub/FunnyGoService/consts"
	"github.com/RustGrub/FunnyGoService/logger"
	"github.com/RustGrub/FunnyGoService/services/FunnyService/models"
	"strconv"
	"time"
)

type RedisGoods struct {
	Cache  cache.Cache
	Logger logger.Logger
}

func NewCacheRepo(c cache.Cache, l logger.Logger) *RedisGoods {
	return &RedisGoods{Cache: c, Logger: l}
}

const goodRedisKey = "funnyservice:good:"
const ttl = time.Minute

// Логи для демонстрации взаимодействия с кэшем, при repo/list ReqID в логах будет nil

func (c *RedisGoods) Set(ctx context.Context, good models.Good) error {
	c.Logger.Warning(fmt.Sprintf("ReqID: %v: trying to set good in cache with name: %v", ctx.Value(consts.ReqID), good.Name))
	key := goodRedisKey + strconv.Itoa(good.GoodID) + ":" + strconv.Itoa(good.ProjectID)
	if err := c.Cache.Set(ctx, key, good, ttl); err != nil {
		return err
	}
	c.Logger.Warning(fmt.Sprintf("ReqID: %v: set good in cache with name: %v", ctx.Value(consts.ReqID), good.Name))
	return nil
}

func (c *RedisGoods) Get(ctx context.Context, goodID, projectID int) *models.Good {
	c.Logger.Warning(fmt.Sprintf("ReqID: %v: trying to get good from cache with id: %v", ctx.Value(consts.ReqID), goodID))
	key := goodRedisKey + strconv.Itoa(goodID) + ":" + strconv.Itoa(projectID)
	value, err := c.Cache.Get(ctx, key)
	if err != nil {
		return nil
	}

	var res models.Good
	if err = json.Unmarshal(value, &res); err != nil {
		return nil
	}
	c.Logger.Warning(fmt.Sprintf("ReqID: %v: got good from cache with name: %v", ctx.Value(consts.ReqID), res.Name))
	return &res
}
