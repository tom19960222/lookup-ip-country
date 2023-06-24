package main

import (
	"net"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	ipRangeList := CreateNetworkDb()

	app.GET("/:ip", func(c *gin.Context) {
		ipStr := c.Param("ip")
		ip := net.ParseIP(ipStr)
		if ip == nil {
			c.String(422, "Invalid IP: %s", ipStr)
			return
		}

		for _, ipRange := range ipRangeList {
			if ipRange.IpRange.Contains(ip) {
				c.String(200, ipRange.CountryCode)
				return
			}
		}
		c.String(404, "IP: %s not in list", ipStr)
	})

	if err := app.Run("0.0.0.0:3000"); err != nil {
		panic(err)
	}

}
