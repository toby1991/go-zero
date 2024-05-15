package svc

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/toby1991/go-zero-utils/bizredis"
	"github.com/toby1991/go-zero-utils/db"
	"github.com/toby1991/go-zero-utils/faktory"
	"github.com/toby1991/go-zero-utils/nsq"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"

    _ent "entgo.io/ent"

    {{.imports}}
)

type ServiceContext struct {
	Config config.Config

	DB       *entClient
	BizRedis bizredis.RedisClient
	Faktory  faktory.FaktoryClient
	Nsq      nsq.NsqClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:c,

		DB:       newEnt(c.Db),
		BizRedis: bizredis.NewRedis(c.BizRedis),
		Faktory:  faktory.NewFaktory(c.Faktory),
		Nsq:      nsq.NewNsq(c.Nsq),
	}
}


// ent client
type entClient struct {
	_conf db.DbConf
	db    *ent.Client
}

func (e *entClient) Db() *ent.Client {
	return e.db
}

func newEnt(conf db.DbConf) *entClient {
	if len(conf.Dsn) <= 0 {
		return nil
	}

	//client, err := ent.Open(conf.DriverName, conf.Dsn)
	drv, err := sql.Open(conf.DriverName, conf.Dsn)
	if err != nil {
		log.Fatalf("failed opening connection to %s: %v", conf.DriverName, err)
		return nil
	}

	// Get the underlying sql.DB object of the driver.
	_db := drv.DB()
	_db.SetMaxIdleConns(10)
	_db.SetMaxOpenConns(100)
	_db.SetConnMaxLifetime(time.Hour)

	var options []ent.Option
	client := ent.NewClient(append(options, ent.Driver(drv))...)

	if conf.Debug {
		client = client.Debug()
	}

	// limit 1000
	//client.Intercept(
	//	intercept.Func(func(ctx context.Context, q intercept.Query) error {
	//		// LimitInterceptor limits the number of records returned from
	//		// the database to 1000, in case Limit was not explicitly set.
	//		if _ent.QueryFromContext(ctx).Limit == nil {
	//			q.Limit(1000)
	//		}
	//		return nil
	//	}),
	//)

	// db migration

	// db migration
	//if err := client.Schema.Create(context.Background(),
	//	migrate.WithForeignKeys(false),
	//	migrate.WithDropIndex(true),
	//	migrate.WithDropColumn(true),
	//); err != nil {
	//	log.Fatalf("failed creating schema resources: %v", err)
	//	return err
	//}
	//// Dump migration changes to stdout.
	if err := client.Schema.WriteTo(context.Background(),
		os.Stdout,
		migrate.WithForeignKeys(false),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		log.Fatalf("failed printing schema changes: %v", err)
	}

	return &entClient{
		_conf: conf,
		db:    client,
	}
}
