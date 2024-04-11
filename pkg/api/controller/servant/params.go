package servant

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPageSize 通用获取页面信息
func GetPageSize(c *gin.Context) (int, int) {
	defaultPageSize := 10
	defaultPageNum := 0

	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	page, errPage := strconv.Atoi(pageStr)
	if errPage != nil {
		page = defaultPageNum
	}

	size, errSize := strconv.Atoi(sizeStr)
	if errSize != nil {
		size = defaultPageSize
	}
	return page, size
}

// GetPathId 通用获取 path 信息
func GetPathId(c *gin.Context) (int, error) {
	return strconv.Atoi(c.Param("id"))
}
