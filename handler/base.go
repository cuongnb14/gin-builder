package handler

import (
	"github.com/cuongnb14/gin-builder/response"
	"github.com/cuongnb14/gin-builder/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FnGetQuery func(c *gin.Context) *gorm.DB

type APIBuilder[M any, S any] struct {
	FnGetQuery FnGetQuery
}

func NewAPIBuilder[M any, S any]() *APIBuilder[M, S] {
	return &APIBuilder[M, S]{}
}

func (b *APIBuilder[M, S]) SetFnGetQuery(fn FnGetQuery) *APIBuilder[M, S] {
	b.FnGetQuery = fn
	return b
}

func (b *APIBuilder[M, S]) BuildListHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		var m []M
		err := b.FnGetQuery(c).Find(&m).Error
		if err != nil {
			response.AbortWithError(c, err)
		}
		s := utils.TranslateList[M, S](&m)
		response.Ok(c, s)
	}
}

func (b *APIBuilder[M, S]) BuildRetrieveHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		pk := c.Param("id")

		var m M
		err := b.FnGetQuery(c).First(&m, pk).Error
		if err != nil {
			response.AbortWithError(c, err)
		}
		s := utils.Translate[S](&m)
		response.Ok(c, s)
	}
}
