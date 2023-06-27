package main

import (
	"log"
	"net/netip"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	ipRangeList := BuildDatabase()

	app.GET("/:ip", func(c *gin.Context) {
		ipStr := c.Param("ip")
		ip, err := netip.ParseAddr(ipStr)
		if err != nil {
			c.String(422, "Invalid IP: %s", ipStr)
			return
		}

		if ip.Is4() {
			for _, ipRange := range ipRangeList.Ipv4 {
				if ipRange.IpRange.Contains(ip) {
					c.String(200, ipRange.CountryCode)
					return
				}
			}
			c.String(404, "IP: %s not in list", ipStr)
			return
		} else if ip.Is6() {
			for _, ipRange := range ipRangeList.Ipv6 {
				if ipRange.IpRange.Contains(ip) {
					c.String(200, ipRange.CountryCode)
					return
				}
			}
			c.String(404, "IP: %s not in list", ipStr)
			return
		} else {
			c.String(422, "Unknown garbage: %s", ipStr)
		}
	})

	log.Println("Starting server ...")
	if err := app.Run("0.0.0.0:3000"); err != nil {
		panic(err)
	}
}
