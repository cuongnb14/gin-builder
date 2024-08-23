package handlerbuilder

import (
	"github.com/cuongnb14/gin-builder/gormfilter"
	"github.com/cuongnb14/gin-builder/pagination"
	"github.com/cuongnb14/gin-builder/response"
	"github.com/cuongnb14/gin-builder/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FnGetQuery func(*gin.Context) *gorm.DB
type FnMapModelToVO[M, S any] func(*gin.Context, M) (S, error)
type FnHook func(*gin.Context) error

type HandlerBuilder[M any, S any] struct {
	FnGetQuery    FnGetQuery
	FilterBuilder *gormfilter.FilterBuilder
	Pagination    *pagination.Pagination
}

func (b *HandlerBuilder[M, S]) SetFnGetQuery(fn FnGetQuery) *HandlerBuilder[M, S] {
	b.FnGetQuery = fn
	return b
}

func (b *HandlerBuilder[M, S]) BuildListHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		var m []M
		query := b.FnGetQuery(c)
		if b.FilterBuilder != nil {
			query = b.FilterBuilder.SetQuery(query).SetRequest(c.Request).GetFilterQuery()
		}
		if b.Pagination != nil {
			page, err := b.Pagination.With(query).Request(c.Request).Response(&m)
			if err != nil {
				response.AbortWithError(c, err)
			}
			pagination.ConvertItems[M, S](page)
			response.OkWithPagination(c, page)
			return
		}

		err := query.Find(&m).Error
		if err != nil {
			response.AbortWithError(c, err)
		}
		s := utils.ConvertList[M, S](&m)
		response.Ok(c, s)

	}
}

func (b *HandlerBuilder[M, S]) BuildRetrieveHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		pk := c.Param("id")

		var m M
		err := b.FnGetQuery(c).First(&m, pk).Error
		if err != nil {
			response.AbortWithError(c, err)
		}
		s := utils.Convert[S](&m)
		response.Ok(c, s)
	}
}
