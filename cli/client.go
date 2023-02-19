package cli

import (
	"github.com/ClickHouse/clickhouse-go"
	"github.com/chuck1024/gd"
	"github.com/chuck1024/gd/databases/mysqldb"
	"github.com/chuck1024/gd/databases/redisdb"
	"github.com/chuck1024/gd/runtime/inject"
	"github.com/jmoiron/sqlx"
	"github.com/leisurelicht/dorm/derror"
)

var Provider *ClientProvider

type ClientProvider struct {
	Mysql      *MysqlClient      `inject:"Mysql"`
	Redis      *RedisClient      `inject:"Redis"`
	ClickHouse *ClickHouseClient `inject:"ClickHose"`
}

func Init(findName string) error {
	f, ok := inject.Find(findName)
	if !ok {
		return derror.DBClientNotFound

	}

	ff, ok := f.(*ClientProvider)
	if !ok {
		return derror.DBClientNotValid
	}

	Provider = ff

	return nil
}

type MysqlClient struct {
	Cli *mysqldb.MysqlClient `inject:"MysqlClient"`
}

func (m *MysqlClient) Start() error {
	return nil
}

type RedisClient struct {
	Cli *redisdb.RedisPoolClient `inject:"RedisClient"`
}

func (r *RedisClient) Start() error {
	return nil
}

type ClickHouseClient struct {
	Cli *sqlx.DB
}

func (ch *ClickHouseClient) Start() (err error) {
	if ch.Cli, err = sqlx.Open("clickhouse", gd.Config("ClickHouse", "addr").String()); err != nil {
		gd.Crash(err)
	}

	if err := ch.Cli.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			gd.Error("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			gd.Error(err)
		}
		return err
	}

	return nil
}
