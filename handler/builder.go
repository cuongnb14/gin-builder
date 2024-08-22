package api_builder

import (
	"fmt"

	"github.com/cuongnb14/gin-builder/response"
	"github.com/cuongnb14/gin-builder/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListAPIBuilder[M any, S any] struct {
	query *gorm.DB
}

func (b *ListAPIBuilder[M, S]) BuildHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		var m []M
		err := b.query.Find(&m).Error
		if err != nil {
			fmt.Println(err)
		}
		s := utils.TranslateList[M, S](&m)
		response.Ok(c, s)
	}
}
