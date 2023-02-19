package clickhouse

import (
	"github.com/leisurelicht/dorm/dao/clickhouse"
	"github.com/leisurelicht/dorm/utils"
	"github.com/leisurelicht/dorm/utils/logger"
	"strings"
)

type ControllerImpl struct {
	dao clickhouse.Impl
}

var _ Controller = (*ControllerImpl)(nil)

func NewCtl(m interface{}) Controller {
	return &ControllerImpl{
		dao: clickhouse.New(m),
	}
}

func (c ControllerImpl) Set(conf map[string]interface{}) Controller {
	c.dao.Set(conf)
	return c
}

func (c ControllerImpl) Filter(filter map[string]interface{}) Controller {
	c.dao.Filter(filter)
	return c
}

func (c ControllerImpl) OrderBy(orderBy string) Controller {
	if orderBy == "" {
		return c
	}

	var orderBySlice []string
	fieldMap := utils.Struct2Map(c.dao.Model(), c.dao.MTag())
	for _, by := range strings.Split(orderBy, ",") {
		by = strings.TrimSpace(by)
		if strings.HasPrefix(by, "-") {
			if _, ok := fieldMap[by[1:]]; ok {
				orderBySlice = append(orderBySlice, by)
			}
		} else {
			if _, ok := fieldMap[by]; ok {
				orderBySlice = append(orderBySlice, by)
			}
		}
	}
	c.dao.OrderBy(orderBySlice)

	return c
}

func (c ControllerImpl) Select(dest interface{}) (err error) {
	logger.PrintCallerInfo()

	if err := c.dao.Select(dest); err != nil {
		return err
	}
	return nil
}
