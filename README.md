# Gin builder

## Package
- `apierror`: Define `APIError` struct that handle by `error_handler` middleware
- `gormfilter`: Auto build gorm query base on gin request param
- `pagination`: Auto paging gorm query
- `handlerbuilder`: Build CURD api base on gorm model


## Usage

```go
package tests

import (
	"fmt"
	"github.com/cuongnb14/gin-builder/gormfilter"
	"github.com/cuongnb14/gin-builder/handlerbuilder"
	"github.com/cuongnb14/gin-builder/pagination"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
	Url filter like: /users?min_age=10&email=user@example.com
*/
func NewUserFilter() *gormfilter.FilterBuilder {
	filterBuilder := &gormfilter.FilterBuilder{}
	filterBuilder.
		AddFilter(gormfilter.Filter[string]{
			Param: "min_age",
			Field: "age",
			Op:    gormfilter.Gte,
		}).
		AddFilter(gormfilter.Filter[string]{
			Param: "email",
			Field: "email",
			Op:    gormfilter.Eq,
		}).
		SetOrdering([]string{"age", "name"})
	return filterBuilder
}

func Test_HandlerBuilder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := GetDB()

	defer func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}()

	CreateUser(db, 5)

	router := gin.Default()

	userHandlerBuilder := handlerbuilder.HandlerBuilder[User, UserVO]{
		FnGetQuery: func(c *gin.Context) *gorm.DB {
			return db.Model(&User{})
		},
		FilterBuilder: NewUserFilter(),
		Pagination:    pagination.NewStandardPagination(),
	}

	router.GET("/users", userHandlerBuilder.BuildListHandler())
	router.GET("/users/:id", userHandlerBuilder.BuildRetrieveHandler())

	t.Run("List users", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users?min_age=44", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		fmt.Println(string(w.Body.Bytes()))
		assert.Contains(t, w.Body.String(), "\"total\":1")
	})

	t.Run("Retrieve user", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		fmt.Println(string(w.Body.Bytes()))
		assert.Contains(t, w.Body.String(), "User 1")
	})
}

```