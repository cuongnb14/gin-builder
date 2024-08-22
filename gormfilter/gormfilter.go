package gormfilter

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GFOperater string

const (
	GFEq  GFOperater = "="
	GFGt  GFOperater = ">"
	GFGte GFOperater = ">="
	GFLt  GFOperater = "<"
	GFLte GFOperater = "<="
)

type IFilter interface {
	GetQuery(value string) (string, interface{}, error)
	GetParam() string
}

type Filter[T any] struct {
	Param string
	Field string
	Op    GFOperater
}

func ConvertString[T any](input string) (T, error) {
	var result T
	var err error

	switch any(result).(type) {
	case int:
		var val int64
		val, err = strconv.ParseInt(input, 10, 0)
		result = any(int(val)).(T)
	case float64:
		var val float64
		val, err = strconv.ParseFloat(input, 64)
		result = any(val).(T)
	case bool:
		var val bool
		val, err = strconv.ParseBool(input)
		result = any(val).(T)
	case uuid.UUID:
		var val uuid.UUID
		val, err = uuid.Parse(input)
		result = any(val).(T)
	case string:
		result = any(input).(T)
	default:
		err = errors.New("unsupported type")
	}

	return result, err
}

func (f Filter[T]) GetParam() string {
	return f.Param
}

func (f Filter[T]) GetQuery(value string) (string, interface{}, error) {
	var where string
	var params interface{}
	params, err := ConvertString[T](value)
	if err != nil {
		return "", nil, err
	}

	where = fmt.Sprintf("%s %s ?", f.Field, f.Op)
	return where, params, nil
}

type FilterBuilder struct {
	Filters  []IFilter
	req      *http.Request
	Ordering []string

	Query *gorm.DB
}

func (b *FilterBuilder) SetQuery(query *gorm.DB) *FilterBuilder {
	b.Query = query
	return b
}

func (b *FilterBuilder) SetOrdering(ordering []string) *FilterBuilder {
	b.Ordering = ordering
	return b
}

func (b *FilterBuilder) SetRequest(req *http.Request) *FilterBuilder {
	b.req = req
	return b
}

func (b *FilterBuilder) AddFilter(filter IFilter) *FilterBuilder {
	b.Filters = append(b.Filters, filter)
	return b
}

func (b *FilterBuilder) BuildFilter() (string, []interface{}, string) {
	var wheres []string
	var params []interface{}
	query := b.req.URL.Query()
	for _, filter := range b.Filters {
		valueParam := query.Get(filter.GetParam())
		if valueParam == "" {
			continue
		}
		w, p, err := filter.GetQuery(valueParam)
		if err != nil {
			slog.Error(err.Error())
		} else {
			wheres = append(wheres, w)
			params = append(params, p)
		}
	}

	ordering := ""

	if len(b.Ordering) > 0 {
		orderingParam := query.Get("sort")

		if orderingParam != "" {
			orderingField := strings.TrimPrefix(orderingParam, "-")

			if slices.Contains(b.Ordering, orderingField) {
				if strings.HasPrefix(orderingParam, "-") {
					ordering = fmt.Sprintf("%s desc", orderingField)
				} else {
					ordering = fmt.Sprintf("%s asc", orderingField)
				}
			}
		}
	}

	return strings.Join(wheres, " and "), params, ordering
}

func (b *FilterBuilder) GetFilterQuery() *gorm.DB {
	wheres, params, ordering := b.BuildFilter()
	if wheres != "" {
		b.Query = b.Query.Where(wheres, params...)
	}
	if ordering != "" {
		b.Query = b.Query.Order(ordering)
	}
	return b.Query
}
