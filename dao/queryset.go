package dao

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

type Operator interface {
	OperatorSQL(operator string) string
}

type QuerySet interface {
	GetFilterSQL() (string, []interface{})
	FilterToSQL(filter map[string]interface{})
	GetOrderBySQL() string
	OrderByToSQL(orderBy []string)
	GetLimitSQL() string
	LimitToSQL(pageSize, pageNum int64)
}

type QuerySetImpl struct {
	filterSql  string
	filterArgs []interface{}
	orderBySql string
	limitSql   string
	Operator
}

var _ QuerySet = (*QuerySetImpl)(nil)

func NewQuerySet(op Operator) *QuerySetImpl {
	return &QuerySetImpl{
		Operator: op,
	}
}

func (p *QuerySetImpl) GetFilterSQL() (string, []interface{}) {
	if p.filterSql != "" {
		return " WHERE" + p.filterSql[3:], p.filterArgs
	}
	return "", p.filterArgs

}

var (
	AND2OR = []string{"AND", "OR"}
	NOT    = []string{"", "NOT"}
)

func (p *QuerySetImpl) FilterToSQL(filter map[string]interface{}) {
	var (
		baseSQL   = " `%s`%s ? "
		fieldName string
		operator  string
		flag      = 0
	)

	if len(filter) == 0 {
		return
	}

	p.filterSql = ""
	p.filterArgs = []interface{}{}

	for fieldLookups, filedValue := range filter {
		p.filterSql += "AND"
		fl := strings.Split(fieldLookups, "__")
		if len(fl) == 1 {
			operator = "exact"
		} else {
			operator = fl[1]
		}
		if len(fl) == 3 && fl[2] == "Q" {
			flag = 1
		}
		fieldName = fl[0]

		op := p.OperatorSQL(operator)
		v := reflect.ValueOf(filedValue)
		switch v.Kind() {
		case reflect.String, reflect.Bool,
			reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Float32, reflect.Float64:
			switch operator {
			case "in", "between":
				log.Panicf("Operator [%s] must be used with slice or array.", operator)
			}
			p.filterSql += fmt.Sprintf(baseSQL, fieldName, op)
			p.filterArgs = append(p.filterArgs, filedValue)
		case reflect.Slice, reflect.Array:
			switch operator {
			case "exact", "exclude", "contains", "icontains":
				p.filterSql += fmt.Sprintf(" ( %s %s ?", fieldName, op) + strings.Repeat(fmt.Sprintf(" %s %s %s ?", AND2OR[flag], fieldName, op), v.Len()-1) + " ) "
			case "in":
				p.filterSql += fmt.Sprintf(" %s %s %s (?"+strings.Repeat(",?", v.Len()-1)+") ", fieldName, NOT[flag], op)
			case "between":
				p.filterSql += fmt.Sprintf(" %s %s %s ? AND ? ", fieldName, NOT[flag], op)
			default:
				continue
			}
			for i := 0; i < v.Len(); i++ {
				p.filterArgs = append(p.filterArgs, v.Index(i).Interface())
			}
		}
	}

	return
}

func (p *QuerySetImpl) GetOrderBySQL() string {
	if strings.HasPrefix(p.orderBySql, ",") {
		return " ORDER BY" + p.orderBySql[1:]
	}
	return ""
}

func (p *QuerySetImpl) OrderByToSQL(orderBy []string) {
	if len(orderBy) <= 0 {
		return
	}
	p.orderBySql = ""
	asc := true
	for _, by := range orderBy {
		p.orderBySql += ","
		by = strings.TrimSpace(by)
		if strings.HasPrefix(by, "-") {
			by = by[1:]
			asc = false
		}

		if asc {
			p.orderBySql += fmt.Sprintf(" %s ASC", by)
		} else {
			p.orderBySql += fmt.Sprintf(" %s DESC", by)
		}
	}
	if strings.HasSuffix(p.orderBySql, ",") {
		p.orderBySql = p.orderBySql[:len(p.orderBySql)-1]
	}

	return
}

func (p *QuerySetImpl) GetLimitSQL() string {
	return p.limitSql
}

func (p *QuerySetImpl) LimitToSQL(pageSize, pageNum int64) {
	p.limitSql = ""
	if pageSize > 0 && pageNum > 0 {
		var offset, limit int64
		offset = (pageNum - 1) * pageSize
		limit = pageSize
		p.limitSql = fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}
	return
}

func stringJoin(elems []string, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("'%s'", elems[0])
	}
	n := len(sep) * (len(elems) - 1)
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString("'")
	b.WriteString(elems[0])
	b.WriteString("'")
	for _, s := range elems[1:] {
		b.WriteString(sep)
		b.WriteString("'")
		b.WriteString(s)
		b.WriteString("'")
	}
	return b.String()
}
