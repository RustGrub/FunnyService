package cache

import (
	"context"
	"github.com/RustGrub/FunnyGoService/config"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	Client *redis.Client
}

func New(cfg *config.Config) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Server + ":" + cfg.Redis.Port, // Адрес и порт Redis сервера
		Password: "",                                      // Пароль, если он установлен
		DB:       cfg.Redis.Database,                      // Номер базы данных
	})
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	return &RedisCache{Client: client}, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	var val string
	var err error
	val, err = r.Client.Get(ctx, key).Result()
	return []byte(val), err
}
