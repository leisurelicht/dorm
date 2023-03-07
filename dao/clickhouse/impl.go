package clickhouse

import (
	"fmt"
	"github.com/leisurelicht/dorm/cli"
	"github.com/leisurelicht/dorm/dao"
	"github.com/leisurelicht/dorm/utils/logger"
	"reflect"
	"strings"
)

var defaultModelTag = "db"

type Impl struct {
	qs   dao.QuerySet
	m    interface{}
	mTag string
}

var _ Dao = (*Impl)(nil)

func New(m interface{}) Impl {
	return Impl{
		qs:   dao.NewQuerySet(newClickHouseOperator()),
		m:    m,
		mTag: defaultModelTag,
	}
}

func (i Impl) Model() interface{} {
	return i.m
}

func (i Impl) MTag() string {
	return i.mTag
}

func (i Impl) Set(conf map[string]interface{}) Dao {
	if model, ok := conf["model"]; ok {
		i.m = model
	}
	if tag, ok := conf["tag"]; ok {
		if t, ok := tag.(string); ok {
			i.mTag = t
		}
	}
	return i
}

func (i Impl) TableName() string {
	return reflect.ValueOf(i.m).MethodByName("TableName").Call([]reflect.Value{})[0].String()
}

func (i Impl) Filter(filter map[string]interface{}) Dao {
	i.qs.FilterToSQL(filter)
	return i
}

func (i Impl) OrderBy(orderBy []string) Dao {
	i.qs.OrderByToSQL(orderBy)
	return i
}

func (i Impl) Limit(pageSize, pageNum int64) Dao {
	i.qs.LimitToSQL(pageSize, pageNum)
	return i
}

func (i Impl) Select(dest interface{}) error {
	logger.PrintCallerInfo()

	sql := fmt.Sprintf("SELECT * FROM %s", i.TableName())

	filterSQL, filterArgs := i.qs.GetFilterSQL()
	sql += filterSQL

	sql += i.qs.GetOrderBySQL()

	sql += i.qs.GetLimitSQL()

	sql = strings.TrimSpace(sql)

	return cli.Provider.ClickHouse.Cli.Select(dest, sql, filterArgs...)
}
