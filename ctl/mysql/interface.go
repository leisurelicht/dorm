package mysql

type Controller interface {
	Model() interface{}
	Set(conf map[string]interface{}) Controller
	Filter(filter map[string]interface{}) Controller
	OrderBy(orderBy string) Controller
	Limit(pageSize, pageNum int64) Controller
	Count() (int64, error)
	Exist() (bool, error)
	Get() (interface{}, error)
	All() ([]interface{}, error)
	List() (int64, []map[string]interface{}, error)
	Create(data map[string]interface{}) error
	Remove(condition map[string]interface{}) (num int64, err error)
	Update(data map[string]interface{}, primaryKeys []string, updateFields []string) error
	Delete(data map[string]interface{}) error
	CreateIfNotExist(data map[string]interface{}, primaryKeys []string) error
	CreateOrUpdate(data map[string]interface{}, primaryKeys []string, updateFields []string) (interface{}, bool, error)
	CreateOrUpdateByModel(primaryKeys []string, updateFields []string) (interface{}, bool, error)
	GetC2CMap(column1, column2 string) (map[string]string, error)
	GetKey2IDMap(key string) (map[string]int64, error)
	GetID2KeyMap(key string) (map[int64]string, error)
	GetID2Map() (map[int64]map[string]interface{}, error)
}
