package filter

import (
	"fmt"
	"strconv"
	"strings"
)

type Filter struct {
	options []Option
}

func NewFilter(options []Option) Filter {
	return Filter{
		options: options,
	}
}

func (f *Filter) String() string {
	builder := strings.Builder{}
	builder.WriteString("[")
	for _, v := range f.options {
		builder.WriteString(v.String() + ", ")
	}
	builder.WriteString("]")
	return builder.String()
}

func (f *Filter) Options() []Option {
	return f.options
}

func (f *Filter) OptionsToSql() string {
	oToSql := func(o Option) string {
		switch o.ParamType {
		case ParamLike:
			return " AND " + o.Name + " LIKE '%" + o.Val + "%'"
		case ParamEq:
			return " AND " + o.Name + " = " + o.Val
		}
		return ""
	}

	sb := strings.Builder{}
	hasPagination := false
	pageSize, pageNum := 1, 1000

	for _, v := range f.options {
		if v.ParamType == ParamPageNum || v.ParamType == ParamPageSize {
			hasPagination = true
			if v.ParamType == ParamPageSize {
				pageSize, _ = strconv.Atoi(v.Val)
			} else {
				pageNum, _ = strconv.Atoi(v.Val)
			}
		} else {
			sb.WriteString(oToSql(v))
		}
	}

	if hasPagination {
		sb.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, (pageNum-1)*pageSize))
	}

	return sb.String()
}
