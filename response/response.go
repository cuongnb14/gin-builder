package response

import (
	"net/http"

	"github.com/cuongnb14/gin-builder/apierror"
	"github.com/cuongnb14/gin-builder/pagination"
	"github.com/gin-gonic/gin"
)

type ResponseOK struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseOKWithPagination struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	//Total   int64       `json:"total"` // TODO: refactor it
	Data interface{} `json:"data"`
}

func Ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, ResponseOK{
		Code:    "Success",
		Message: "Success",
		Data:    data,
	})
}

func OkWithPagination(c *gin.Context, page *pagination.Page) {
	c.JSON(http.StatusOK, ResponseOKWithPagination{
		Code:    "Success",
		Message: "Success",
		Data:    page.Items,
		//Total:   page.Total,
	})
}

func AbortWithAPIError(ctx *gin.Context, err *apierror.APIError) {
	_ = ctx.Error(err).SetType(gin.ErrorTypePublic)
	ctx.Abort()
}

func AbortWithError(ctx *gin.Context, err error) {
	_ = ctx.Error(err).SetType(gin.ErrorTypePublic)
	ctx.Abort()
}

func AbortWithBindError(ctx *gin.Context, err error) {
	_ = ctx.Error(err).SetType(gin.ErrorTypeBind)
	ctx.Abort()
}
