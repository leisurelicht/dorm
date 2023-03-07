package mysql

type Dao interface {
	Model() interface{}
	MTag() string
	Set(conf map[string]interface{}) Dao
	Filter(filter map[string]interface{}) Dao
	OrderBy(orderBy []string) Dao
	Limit(pageSize, pageNum int64) Dao
	Create() error
	Update(primaryKeys []string, updateFields []string) error
	CreateOrUpdate(primaryKeys []string, updateFields []string) (err error)
	Count() (int64, error)
	Get() (interface{}, error)
	All() ([]interface{}, error)
}
