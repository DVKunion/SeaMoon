package control

import (
	"github.com/gin-gonic/gin"
)

type StatisticController struct {
}

func (sc StatisticController) Get(c *gin.Context) {
	// 目前数据结构对统计不太友好，放到下一期做大盘吧
}
