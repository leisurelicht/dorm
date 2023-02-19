package clickhouse

var operators = map[string]string{
	"exact":       "=",
	"exclude":     "!=",
	"iexact":      "LIKE",
	"contains":    "LIKE",
	"icontains":   "ILIKE",
	"gt":          ">",
	"gte":         ">=",
	"lt":          "<",
	"lte":         "<=",
	"startswith":  "LIKE",
	"endswith":    "LIKE",
	"istartswith": "ILIKE",
	"iendswith":   "ILIKE",
	"in":          "IN",
}

type clickHouseOperator struct{}

func newClickHouseOperator() *clickHouseOperator {
	return &clickHouseOperator{}
}

func (d *clickHouseOperator) OperatorSQL(operator string) string {
	return operators[operator]
}
