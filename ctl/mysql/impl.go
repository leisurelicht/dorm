package mysql

import (
	"errors"
	"github.com/leisurelicht/dorm/dao/mysql"
	"github.com/leisurelicht/dorm/derror"
	"github.com/leisurelicht/dorm/utils"
	"github.com/leisurelicht/dorm/utils/logger"
	"strings"
)

type Impl struct {
	dao mysql.Dao
}

var _ Controller = (*Impl)(nil)

func NewCtl(model interface{}) Controller {
	return Impl{
		dao: mysql.New(model),
	}
}

func (i Impl) Model() interface{} {
	return i.dao.Model()
}

func (i Impl) Set(conf map[string]interface{}) Controller {
	i.dao.Set(conf)
	return i
}

func (i Impl) Filter(filter map[string]interface{}) Controller {
	i.dao.Filter(filter)
	return i
}

func (i Impl) OrderBy(orderBy string) Controller {
	if orderBy == "" {
		return i
	}

	var orderBySlice []string
	fieldMap := utils.Struct2Map(i.dao.Model(), i.dao.MTag())
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
	i.dao.OrderBy(orderBySlice)
	return i
}

func (i Impl) Limit(pageSize, pageNum int64) Controller {
	i.dao.Limit(pageSize, pageNum)
	return i
}

func (i Impl) Count() (int64, error) {
	logger.PrintCallerInfo()

	return i.dao.Count()
}

func (i Impl) Exist() (bool, error) {
	logger.PrintCallerInfo()

	if num, err := i.dao.Count(); err != nil {
		return false, err
	} else if num > 0 {
		return true, nil
	}

	return false, nil
}

func (i Impl) Get() (interface{}, error) {
	logger.PrintCallerInfo()

	return i.dao.Get()
}

func (i Impl) All() (res []interface{}, err error) {
	logger.PrintCallerInfo()

	return i.dao.All()
}

func (i Impl) List() (total int64, res []map[string]interface{}, err error) {
	logger.PrintCallerInfo()

	if total, err = i.dao.Count(); err != nil {
		return
	}

	if data, err := i.dao.All(); err != nil {
		return 0, nil, err
	} else {
		for _, d := range data {
			res = append(res, utils.Struct2Map(d, i.dao.MTag()))
		}
		return total, res, nil
	}
}

func (i Impl) Create(data map[string]interface{}) (err error) {
	logger.PrintCallerInfo()

	if err = utils.DecodeByTag(data, i.dao.Model(), i.dao.MTag()); err != nil {
		return err
	}

	if err = i.dao.Create(); err != nil {
		return err
	}

	return nil
}

func (i Impl) Update(
	data map[string]interface{},
	primaryKeys []string,
	updateFields []string,
) (err error) {
	logger.PrintCallerInfo()

	if err = utils.DecodeByTag(data, i.dao.Model(), i.dao.MTag()); err != nil {
		return err
	}

	if err = i.dao.Update(primaryKeys, updateFields); err != nil {
		return err
	}

	return nil
}

func (i Impl) Delete(data map[string]interface{}) (err error) {
	logger.PrintCallerInfo()

	primaryKeys, _, err := utils.Map2SliceE(data)
	if err != nil {
		return err
	}
	updateFields := []string{"is_deleted"}

	data["is_deleted"] = true

	if err = utils.DecodeByTag(data, i.dao.Model(), i.dao.MTag()); err != nil {
		return err
	}

	if err = i.dao.Update(primaryKeys.([]string), updateFields); err != nil {
		return err
	}

	return nil
}

func (i Impl) CreateIfNotExist(
	data map[string]interface{}, primaryKeys []string,
) (err error) {
	logger.PrintCallerInfo()

	filter := make(map[string]interface{})

	for _, key := range primaryKeys {
		filter[key] = data[key]
	}

	if exist, err := i.Filter(filter).Exist(); err != nil {
		return err
	} else if exist {
		return nil
	}

	if err = utils.DecodeByTag(data, i.dao.Model(), i.dao.MTag()); err != nil {
		return err
	}

	if err := i.dao.Create(); err != nil {
		return err
	}

	return nil
}

func (i Impl) CreateOrUpdate(
	data map[string]interface{}, primaryKeys []string, updateFields []string,
) (obj interface{}, created bool, err error) {
	logger.PrintCallerInfo()

	filter := make(map[string]interface{})
	for _, key := range primaryKeys {
		filter[key] = data[key]
	}

	if obj, err = i.dao.Filter(filter).Get(); err != nil {
		if !errors.Is(err, derror.DoesNotExist) {
			return nil, false, err
		}
		if err := i.Create(data); err != nil {
			return nil, false, err
		}
		created = true
	} else {
		objData := utils.Struct2Map(obj, i.dao.MTag())
		for _, key := range updateFields {
			if v, ok := objData[key]; !ok {
				continue
			} else {
				if v == data[key] {
					continue
				}
				if err = i.Update(data, primaryKeys, updateFields); err != nil {
					return nil, false, err
				}
				break
			}
		}
	}

	if obj, err = i.dao.Filter(filter).Get(); err != nil {
		return nil, false, err
	}

	return obj, created, nil
}

func (i Impl) CreateOrUpdateByModel(
	primaryKeys []string, updateFields []string,
) (obj interface{}, created bool, err error) {
	logger.PrintCallerInfo()

	filter := utils.Struct2MapFilterByKeys(i.dao.Model(), i.dao.MTag(), primaryKeys)

	if obj, err = i.dao.Filter(filter).Get(); err != nil {
		if !errors.Is(err, derror.DoesNotExist) {
			return nil, false, err
		}
		if err := i.dao.Create(); err != nil {
			return nil, false, err
		}
		created = true
	} else {
		objData := utils.Struct2Map(obj, i.dao.MTag())
		modelData := utils.Struct2Map(i.dao.Model(), i.dao.MTag())
		for _, key := range updateFields {
			if objData[key] == modelData[key] {
				continue
			}
			if err = i.dao.Update(primaryKeys, updateFields); err != nil {
				return nil, false, err
			}
			break
		}
	}

	if obj, err = i.dao.Filter(filter).Get(); err != nil {
		return nil, false, err
	}

	return obj, created, nil
}

func (i Impl) GetC2CMap(column1, column2 string) (res map[string]string, err error) {
	logger.PrintCallerInfo()

	if data, err := i.dao.All(); err != nil {
		return nil, err
	} else {
		res = make(map[string]string, len(data))

		for _, d := range data {
			column1Value := utils.GetValueByTag(d, i.dao.MTag(), column1).(string)
			column2Value := utils.GetValueByTag(d, i.dao.MTag(), column2).(string)
			res[column1Value] = column2Value
		}
	}
	return res, nil
}

func (i Impl) GetKey2IDMap(key string) (res map[string]int64, err error) {
	logger.PrintCallerInfo()

	if data, err := i.dao.All(); err != nil {
		return nil, err
	} else {
		res = make(map[string]int64, len(data))

		for _, d := range data {
			name := utils.GetValueByTag(d, i.dao.MTag(), key).(string)
			id := utils.GetValueByTag(d, i.dao.MTag(), "id").(int64)
			res[name] = id
		}
	}
	return res, nil
}

func (i Impl) GetID2KeyMap(key string) (res map[int64]string, err error) {
	logger.PrintCallerInfo()

	if data, err := i.dao.All(); err != nil {
		return nil, err
	} else {
		res = make(map[int64]string, len(data))

		for _, d := range data {
			id := utils.GetValueByTag(d, i.dao.MTag(), "id").(int64)
			name := utils.GetValueByTag(d, i.dao.MTag(), key).(string)
			res[id] = name
		}
	}
	return res, nil
}

func (i Impl) GetID2Map() (res map[int64]map[string]interface{}, err error) {
	logger.PrintCallerInfo()

	if data, err := i.dao.All(); err != nil {
		return nil, err
	} else {
		res = make(map[int64]map[string]interface{}, len(data))

		for _, d := range data {
			id := utils.GetValueByTag(d, i.dao.MTag(), "id").(int64)
			res[id] = utils.Struct2Map(d, i.dao.MTag())
		}
	}
	return res, nil
}
