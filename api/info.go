package api

import (
	"github.com/gin-gonic/gin"
)

func GetInfo(c *gin.Context) {

	status, err := Runtime(c).GetInfo()
	if err != nil {
		c.String(500, err.Error())
	} else {
		c.JSON(200, status)
	}
}
