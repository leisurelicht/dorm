package clickhouse

type Controller interface {
	Set(conf map[string]interface{}) Controller
	Filter(filter map[string]interface{}) Controller
	OrderBy(orderBy string) Controller
	Select(dest interface{}) error
}
