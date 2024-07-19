package rest

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/api/constant"
	"github.com/sagernet/sing-box/api/rest/rp"
	rq "github.com/sagernet/sing-box/api/rest/rq"
	"github.com/sagernet/sing-box/experimental/v2rayapi"
	"net/http"
)

func getAUserStat(uuid string) (rp.StatRp, error) {
	response, err := box.BoxInstance.
		Router().
		V2RayServer().
		StatsService().(v2rayapi.StatsServiceServer).
		QueryStats(
			context.Background(),
			&v2rayapi.QueryStatsRequest{Regexp: true,
				Patterns: []string{"^user>>>" + uuid + ">>>traffic>>>uplink$", "^user>>>" + uuid + ">>>traffic>>>downlink$"}})

	if err != nil {
		return rp.StatRp{}, err
	} else {
		if response.Stat != nil {
			if len(response.Stat) == 2 {
				return rp.StatRp{
					Id:       uuid,
					Uplink:   response.Stat[0].Value,
					Downlink: response.Stat[1].Value,
				}, nil
			} else if len(response.Stat) > 2 {
				return rp.StatRp{}, errors.New("Found more than one in v2ay api with uuid " + uuid)
			} else if len(response.Stat) < 2 {
				return rp.StatRp{}, errors.New("User " + uuid + " not found or inactivate in v2ray api.")
			}
		} else {
			return rp.StatRp{}, errors.New("User " + uuid + " not found or inactivate in v2ray api.")
		}
	}
	return rp.StatRp{}, errors.New("nothing")
}

func GetAllUsersStats(c *gin.Context) {
	if box.BoxInstance.Router().V2RayServer() != nil {
		var req rq.GetStatsRq
		if err := c.ShouldBindJSON(&req); err != nil {
			req.Reset = true
		}
		listUserStats := make([]rp.StatRp, 0)
		for uuid, _ := range constant.InRamUsersUUID {
			if req.Reset {
				stat, err := getAUserStat(uuid)
				if err == nil {
					listUserStats = append(listUserStats, stat)
				}
				box.BoxInstance.Router().V2RayServer().StatsService().EditUser(uuid, true)
				box.BoxInstance.Router().V2RayServer().StatsService().EditUser(uuid, false)
			} else {
				stat, err := getAUserStat(uuid)
				if err == nil {
					listUserStats = append(listUserStats, stat)
				}
			}
		}
		if len(listUserStats) > 0 {
			c.JSON(http.StatusOK, rp.StatRpList{ListUserStats: listUserStats})
		} else if len(listUserStats) == 0 && !req.Reset {
			c.JSON(http.StatusOK, gin.H{"error": "No active user found."})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "V2ray api stats not initialized."})
	}
}
func GetAUserStat(c *gin.Context) {
	if box.BoxInstance.Router().V2RayServer() != nil {
		var req rq.GetAUserStatRq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		stat, err := getAUserStat(req.UUID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, stat)
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "V2ray api stats not initialized."})
	}
}
