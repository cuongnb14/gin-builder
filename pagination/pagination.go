package pagination

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/cuongnb14/gin-builder/utils"
	"gorm.io/gorm"
)

func NewStandardPagination() *Pagination {
	return &Pagination{
		defaultSize: 24,
		maxSize:     100,
		withCount:   true,
	}
}

type Pagination struct {
	defaultSize   int
	maxSize       int
	Query         *gorm.DB
	RequestParams *RequestParams
	withCount     bool
}

type Page struct {
	Items interface{}
	Total *int64
}

type RequestParams struct {
	Size int
	Page int
}

func (p *Pagination) With(query *gorm.DB) *Pagination {
	p.Query = query
	return p
}

func (p *Pagination) SetCount(enableCount bool) *Pagination {
	p.withCount = enableCount
	return p
}

func (p *Pagination) Request(request *http.Request) *Pagination {
	query := request.URL.Query()
	requestParams := &RequestParams{}
	if i, e := strconv.Atoi(query.Get("size")); nil == e {
		requestParams.Size = slices.Min([]int{i, p.maxSize})
	} else {
		requestParams.Size = p.defaultSize
	}

	if i, e := strconv.Atoi(query.Get("page")); nil == e {
		requestParams.Page = i
	} else {
		requestParams.Page = 1
	}
	p.RequestParams = requestParams
	return p
}

func (p *Pagination) Response(results interface{}) (*Page, error) {
	var total *int64
	dbs := p.Query.Statement.DB.Session(&gorm.Session{NewDB: true})
	query := dbs.Unscoped().Table("(?) AS s", p.Query)
	if p.withCount {
		var zero int64
		total = &zero
		query = query.Count(total).Limit(p.RequestParams.Size).Offset((p.RequestParams.Page - 1) * p.RequestParams.Size)
	} else {
		query = query.Limit(p.RequestParams.Size).Offset((p.RequestParams.Page - 1) * p.RequestParams.Size)
	}
	err := query.Find(results).Error
	if err != nil {
		return nil, err
	}

	return &Page{
		Items: results,
		Total: total,
	}, nil
}

func ConvertItems[F any, T any](page *Page) {
	newType, _ := page.Items.(*[]F)
	page.Items = utils.ConvertList[F, T](newType)
}
