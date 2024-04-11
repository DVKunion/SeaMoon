package servant

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/DVKunion/SeaMoon/pkg/system/errors"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

func RawMsg(c *gin.Context, name string, text []byte) {
	c.Header("Content-Disposition", "attachment; filename="+name)
	c.Data(200, "application/octet-stream; charset=utf-8", text)
	return
}

// SuccessMsg 通用正常响应
func SuccessMsg(c *gin.Context, total int64, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"total":   total,
		"data":    data,
	})
	return
}

// ErrorMsg 通用错误响应
func ErrorMsg(c *gin.Context, code int, err error) {
	if err == nil {
		err = errors.ApiError(xlog.ApiCommonError, nil)
	}
	c.JSON(code, gin.H{
		"success": false,
		"code":    code,
		"msg":     err.Error(),
	})
}
