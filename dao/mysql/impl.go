package mysql

import (
	"fmt"
	"github.com/leisurelicht/dorm/cli"
	"github.com/leisurelicht/dorm/dao"
	"github.com/leisurelicht/dorm/derror"
	"github.com/leisurelicht/dorm/utils"
	"github.com/leisurelicht/dorm/utils/logger"
	"log"
	"reflect"
	"strings"
	"time"
)

var defaultModelTag = "mysqlField"

type Impl struct {
	qs   dao.QuerySet
	m    interface{}
	mTag string
}

var _ Dao = (*Impl)(nil)

func New(m interface{}) Dao {
	if m == nil {
		log.Panicln("New Mysql Dao Error: model is nil")
		return nil
	}

	t := reflect.TypeOf(m)
	if t.Kind() != reflect.Ptr {
		log.Panicf("New Mysql Dao Error: model [%+v] is not a pointer", t.Name())
		return nil
	}

	return Impl{qs: dao.NewQuerySet(newMysqlOperator()), m: m, mTag: defaultModelTag}
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

func (i Impl) Limit(pageSize, pageNum int) Dao {
	i.qs.LimitToSQL(pageSize, pageNum)
	return i
}

func (i Impl) Count() (num int64, err error) {
	logger.PrintCallerInfo()

	sql := fmt.Sprintf("SELECT count(*) FROM %s", i.TableName())

	filterSQL, filterArgs := i.qs.GetFilterSQL()
	sql += filterSQL

	sql = strings.TrimSpace(sql)

	return cli.Provider.Mysql.Cli.GetCount(sql, filterArgs...)
}

func (i Impl) Get() (data interface{}, err error) {
	logger.PrintCallerInfo()

	sql := fmt.Sprintf("SELECT ? FROM %s", i.TableName())

	filterSQL, filterArgs := i.qs.GetFilterSQL()
	sql += filterSQL

	sql += i.qs.GetOrderBySQL()

	sql += i.qs.GetLimitSQL()

	sql = strings.TrimSpace(sql)

	if data, err = cli.Provider.Mysql.Cli.Query(i.Model(), sql, filterArgs...); err != nil {
		return nil, err
	}

	if data == nil {
		return nil, derror.DoesNotExist
	}

	return data, nil
}

func (i Impl) All() (data []interface{}, err error) {
	logger.PrintCallerInfo()

	sql := fmt.Sprintf("SELECT ? FROM %s", i.TableName())

	filterSQL, filterArgs := i.qs.GetFilterSQL()
	sql += filterSQL

	sql += i.qs.GetOrderBySQL()

	sql += i.qs.GetLimitSQL()

	sql = strings.TrimSpace(sql)

	return cli.Provider.Mysql.Cli.QueryList(i.Model(), sql, filterArgs...)
}

func (i Impl) Create() (err error) {
	logger.PrintCallerInfo()

	// 创建时自动给 created 及 updated 赋值
	utils.SetValuesByTag(
		i.Model(), i.MTag(),
		map[string]interface{}{
			"created": time.Now().Unix(),
			"updated": time.Now().Unix(),
		},
	)

	return cli.Provider.Mysql.Cli.Add(i.TableName(), i.Model(), false)
}

func (i Impl) Update(primaryKeys []string, updateFields []string) (err error) {
	logger.PrintCallerInfo()

	var (
		pksMap = make(map[string]interface{})
		upfMap = make(map[string]interface{})
		upf    []string
	)

	// 更新时自动更新 updated
	utils.SetValuesByTag(
		i.Model(), i.MTag(),
		map[string]interface{}{
			"updated": time.Now().Unix(),
		},
	)
	updateFields = append(updateFields, "updated")
	updateFields = utils.DuplicateString(updateFields)

	pksMap = utils.StringSlice2Map(primaryKeys)
	upfMap = utils.StringSlice2Map(updateFields)

	t := reflect.TypeOf(i.m).Elem()
	v := reflect.ValueOf(i.m).Elem()

	for k := 0; k < t.NumField(); k++ {
		key := t.Field(k).Tag.Get(i.MTag())
		if _, ok := pksMap[key]; ok {
			pksMap[key] = v.Field(k).Interface()
		}

		if _, ok := upfMap[key]; ok {
			upf = append(upf, key)
		}
	}

	return cli.Provider.Mysql.Cli.Update(i.TableName(), i.Model(), pksMap, upf)
}

func (i Impl) CreateOrUpdate(primaryKeys []string, updateFields []string) (err error) {
	logger.PrintCallerInfo()

	utils.SetValuesByTag(
		i.Model(), i.MTag(),
		map[string]interface{}{
			"created": time.Now().Unix(),
			"updated": time.Now().Unix(),
		},
	)
	updateFields = append(updateFields, "updated")
	updateFields = utils.DuplicateString(updateFields)

	_, err = cli.Provider.Mysql.Cli.InsertOrUpdateOnDup(
		i.TableName(), i.Model(), primaryKeys, updateFields, true,
	)

	return err
}
