package clickhouse

type Dao interface {
	Model() interface{}
	MTag() string
	Set(conf map[string]interface{}) Dao
	Filter(filter map[string]interface{}) Dao
	OrderBy(orderBy []string) Dao
	Limit(pageSize, pageNum int) Dao
	Select(dest interface{}) error
}
