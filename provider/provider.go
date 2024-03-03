package provider

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/RustGrub/FunnyGoService/broker"
	"github.com/RustGrub/FunnyGoService/cache"
	"github.com/RustGrub/FunnyGoService/config"
	"github.com/RustGrub/FunnyGoService/logger"
	"github.com/RustGrub/FunnyGoService/sql/clhouse"
	"net/url"
	// Охота "симметричность" в двух разных проектах, но не знаю, как это грамотно сделать...
	fsgoods "github.com/RustGrub/FunnyGoService/services/FunnyService/repo/goods"
	lsgoods "github.com/RustGrub/FunnyGoService/services/ListeningService/repo/goods"
	"github.com/RustGrub/FunnyGoService/sql/txmanager"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Provider interface {
	GetTxManager() txmanager.TxManager
	GetGoodsRepository() fsgoods.FsRepository
	GetGoodsCache() fsgoods.Cache
	GetBroker() broker.MessageBroker
	GetLsRepo() lsgoods.Repository
}

type provider struct {
	TxManager txmanager.TxManager

	FsGoodsRepo  fsgoods.FsRepository
	FsGoodsCache fsgoods.Cache

	LsRepo lsgoods.Repository

	Broker broker.MessageBroker
}

func New(cfg *config.Config, logger logger.Logger) Provider {
	// DB for FunnyService (pg)
	pool, err := createPgxPool(cfg.ProjectsDb)
	if err != nil {
		logger.Fatal(fmt.Errorf("error in provider creating pgx pool: %v", err))
	}
	fsTxMng := txmanager.New(pool)
	fsGoodsRepo := fsgoods.NewPgRepo(pool, logger)

	// DB for listeningService (clickhouse)
	ch, err := createCHouseConn(cfg.House)
	if err != nil {
		logger.Fatal(fmt.Errorf("error in provider creating cHouseConn: %v", err))
	}
	lsRepo := lsgoods.New(clhouse.New(ch), logger)

	// Cache for FunnyService
	redisCache, err := cache.New(cfg)
	if err != nil {
		logger.Fatal(fmt.Errorf("error in provider creating redisCache: %v", err))
	}

	goodsCache := fsgoods.NewCacheRepo(redisCache, logger)

	// Msg brokers for both
	nts, err := broker.New(cfg)
	if err != nil {
		logger.Fatal(fmt.Errorf("error in provider creating nats con: %v", err))
	}

	return &provider{TxManager: fsTxMng, FsGoodsRepo: fsGoodsRepo, FsGoodsCache: goodsCache, Broker: nts, LsRepo: lsRepo}
}

func (p *provider) GetTxManager() txmanager.TxManager {
	return p.TxManager
}
func (p *provider) GetGoodsRepository() fsgoods.FsRepository {
	return p.FsGoodsRepo
}
func (p *provider) GetGoodsCache() fsgoods.Cache {
	return p.FsGoodsCache
}
func (p *provider) GetBroker() broker.MessageBroker {
	return p.Broker
}
func (p *provider) GetLsRepo() lsgoods.Repository {
	return p.LsRepo
}

func createPgxPool(cfg *config.DatabaseCfg) (*pgxpool.Pool, error) {
	query := url.Values{}
	query.Add("dbname", cfg.Database)
	query.Add("sslmode", "disable")

	host := cfg.Server + ":" + cfg.Port

	dbUrl := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.Username, cfg.Password),
		Host:     host,
		RawQuery: query.Encode(),
	}

	pool, err := pgxpool.New(context.Background(), dbUrl.String())
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func createCHouseConn(cfg *config.DatabaseCfg) (*sql.DB, error) {
	conn, err := sql.Open("chhttp", "http://user:password@127.0.0.1:8123/default")

	if err != nil {
		return nil, err
	}

	fmt.Println("123")
	fmt.Println(conn.Ping())
	fmt.Println("1234")
	return nil, nil
}
