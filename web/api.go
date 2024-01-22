package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"proxy/data"
)

func InitApi(host string, port int) {
	c := gin.Default()

	c.GET("/count", func(c *gin.Context) {
		dataCount, _ := data.GlobalEngine.CountProxiesByTypeAndCountry()
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "success",
			"data": dataCount,
		})
	})

	c.GET("/list", func(c *gin.Context) {
		isJson := c.DefaultQuery("json", "")
		country := c.DefaultQuery("country", "")
		_type := c.DefaultQuery("type", "")

		dataProxyList, _ := data.GlobalEngine.GetProxyList(country, _type)
		if isJson != "" {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"msg":  "success",
				"data": dataProxyList,
			})
		} else {
			textList := ""
			for _, v := range dataProxyList {
				textList += v.Proxy + "\n"
			}
			c.String(http.StatusOK, textList)
		}

	})

	c.Run(fmt.Sprintf("%s:%d", host, port))
}
