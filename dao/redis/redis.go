package redis

import (
	"fmt"
	"github.com/leisurelicht/dorm/cli"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

type Int64Dao interface {
	Get() (int64, error)
	Set(int64) error
}
type Int64DaoImpl struct {
	Key string
}

func NewInt64Dao(key string) Int64Dao {
	return &Int64DaoImpl{Key: key}
}
func (i *Int64DaoImpl) Set(value int64) error {
	return cli.Provider.Redis.Cli.Set(i.Key, strconv.FormatInt(value, 10))
}
func (i *Int64DaoImpl) Get() (int64, error) {
	res, err := cli.Provider.Redis.Cli.Get(i.Key)
	if err != nil && err == redis.ErrNil {
		return 0, nil
	}
	return strconv.ParseInt(res, 10, 64)
}

type StringDao interface {
	Set(string) error
	Get() string
}
type StringDaoImpl struct {
	Key string
}

func NewStringDao(key string) StringDao {
	return &StringDaoImpl{Key: key}
}
func (i *StringDaoImpl) Set(value string) error {
	if _, err := cli.Provider.Redis.Cli.RPush(i.Key, value); err != nil {
		return err
	}
	return nil
}

func (i *StringDaoImpl) Get() string {
	res, err := cli.Provider.Redis.Cli.Get(i.Key)
	if err != nil && err == redis.ErrNil {
		return ""
	}
	return res
}

type SetDao interface {
	Set([]string) error
}
type SetDaoImpl struct {
	Key string
}

func NewSetDao(key string) SetDao {
	return &SetDaoImpl{Key: key}
}
func (i *SetDaoImpl) Set(value []string) error {
	return cli.Provider.Redis.Cli.SAdd(i.Key, value)
}

type ByteListDao interface {
	Len() (int64, error)
	Set([]byte) error
	List() ([][]byte, error)
}

type ByteListDaoImpl struct {
	Key string
}

func NewByteListDao(key string) ByteListDao {
	return &ByteListDaoImpl{Key: key}
}

func (b *ByteListDaoImpl) Len() (int64, error) {
	return redis.Int64(cli.Provider.Redis.Cli.Do("LLEN", b.Key))
}

func (b *ByteListDaoImpl) Set(value []byte) error {
	_, err := redis.Int64(cli.Provider.Redis.Cli.Do("RPUSH", b.Key, value))
	return err
}

func (b *ByteListDaoImpl) List() ([][]byte, error) {
	stop, err := b.Len()
	if err != nil {
		return nil, fmt.Errorf("get byte list [%s] length error: %s", b.Key, err)
	}
	res, err := redis.ByteSlices(cli.Provider.Redis.Cli.Do("LRANGE", b.Key, 0, stop))
	if err != nil {
		return nil, fmt.Errorf("lrange byte list [%s] error: %s", b.Key, err)
	}
	if err == redis.ErrNil {
		return [][]byte{}, nil
	}

	err = cli.Provider.Redis.Cli.Del(b.Key)

	return res, err
}
