package mysql

var operators = map[string]string{
	"exact":       "=",
	"exclude":     "!=",
	"iexact":      "LIKE",
	"contains":    "LIKE BINARY",
	"icontains":   "LIKE",
	"gt":          ">",
	"gte":         ">=",
	"lt":          "<",
	"lte":         "<=",
	"startswith":  "LIKE BINARY",
	"endswith":    "LIKE BINARY",
	"istartswith": "LIKE",
	"iendswith":   "LIKE",
	"in":          "IN",
	"between":     "BETWEEN",
}

type mysqlOperator struct{}

func newMysqlOperator() *mysqlOperator {
	return &mysqlOperator{}
}

func (d *mysqlOperator) OperatorSQL(operator string) string {
	return operators[operator]
}
