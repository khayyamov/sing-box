package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	box "github.com/sagernet/sing-box"
	rq "github.com/sagernet/sing-box/api/rest/rq"
	"github.com/sagernet/sing-box/experimental/v2rayapi"
	"net/http"
)

func getUserFullUsage(c *gin.Context) {
	if box.BoxInstance.Router().V2RayServer() != nil {
		var req rq.GetUserFullUsageRq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response, err := box.BoxInstance.
			Router().
			V2RayServer().
			StatsService().(v2rayapi.StatsServiceServer).
			QueryStats(
				context.Background(),
				&v2rayapi.QueryStatsRequest{Regexp: true,
					Patterns: []string{"^user>>>" + req.UUID + ">>>traffic>>>uplink$", "^user>>>" + req.UUID + ">>>traffic>>>downlink$"}})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			if response.Stat != nil {
				if len(response.Stat) == 2 {
					count := response.Stat[0].Value + response.Stat[0].Value
					c.JSON(http.StatusOK, gin.H{"full_usage": count})
				} else if len(response.Stat) > 2 {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Multiple user found with regex: " + "^user>>>" + req.UUID + ">>>traffic>>>$"})
				} else if len(response.Stat) < 2 {
					c.JSON(http.StatusBadRequest, gin.H{"error": "No user found with regex: " + "^user>>>" + req.UUID + ">>>traffic>>>$"})
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "No user found with regex: " + "^user>>>" + req.UUID + ">>>traffic>>>$"})
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "V2ray api stats not initialized."})
	}
}
