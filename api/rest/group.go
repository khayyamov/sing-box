package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	global_import "github.com/sagernet/sing-box/api/rest/rq"
	"github.com/sagernet/sing-box/option"
	shadowtls "github.com/sagernet/sing-shadowtls"
	"github.com/sagernet/sing/common/auth"
	"io"
	"net/http"
)

func AddToAll(c *gin.Context) {
	domesticLogicGroup(c, false)
}
func DeleteToAll(c *gin.Context) {
	domesticLogicGroup(c, true)
}

func domesticLogicGroup(c *gin.Context, delete bool) {
	var rq global_import.GlobalModel
	var rqArr []global_import.GlobalModel

	bodyAsByteArray, _ := io.ReadAll(c.Request.Body)
	jsonBody := string(bodyAsByteArray)
	haveErr := false

	err := json.Unmarshal([]byte(jsonBody), &rq)
	if err == nil {
		haveErr = false
		EditGroupUsers(c, []global_import.GlobalModel{rq}, delete)
	} else {
		haveErr = true
	}
	err = json.Unmarshal([]byte(jsonBody), &rqArr)
	if err == nil {
		haveErr = false
		EditGroupUsers(c, rqArr, delete)
	} else {
		haveErr = true
	}

	if haveErr {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func EditGroupUsers(c *gin.Context, newUsers []global_import.GlobalModel, delete bool) {
	hysteria := make([]option.HysteriaUser, len(newUsers))
	hysteria2 := make([]option.Hysteria2User, len(newUsers))
	naive := make([]auth.User, len(newUsers))
	shadowsocks_multi := make([]option.ShadowsocksUser, len(newUsers))
	shadowsocks_relay := make([]option.ShadowsocksDestination, len(newUsers))
	shadowtlsArr := make([]shadowtls.User, len(newUsers))
	trojan := make([]option.TrojanUser, len(newUsers))
	tuic := make([]option.TUICUser, len(newUsers))
	vless := make([]option.VLESSUser, len(newUsers))
	vmess := make([]option.VMessUser, len(newUsers))
	for i := range newUsers {
		if newUsers[i].AddToAll || newUsers[i].Hysteria {
			hysteria = append(hysteria, option.HysteriaUser{
				Name:       newUsers[i].Name,
				Auth:       []byte(newUsers[i].Auth),
				AuthString: newUsers[i].AuthString,
			})
		}

		if newUsers[i].AddToAll || newUsers[i].Hysteria2 {
			hysteria2 = append(hysteria2, option.Hysteria2User{
				Name:     newUsers[i].Name,
				Password: newUsers[i].Name,
			})
			EditHysteria2Users(c, []option.Hysteria2User{{
				Name:     newUsers[i].Name,
				Password: newUsers[i].Name,
			}}, delete)
		}

		if newUsers[i].AddToAll || newUsers[i].Naive {
			naive = append(naive, auth.User{
				Username: newUsers[i].Username,
				Password: newUsers[i].Name,
			})
			EditNaiveUsers(c, []auth.User{{
				Username: newUsers[i].Username,
				Password: newUsers[i].Name,
			}}, delete)
		}

		if newUsers[i].AddToAll || newUsers[i].Shadowsocks_multi {
			shadowsocks_multi = append(shadowsocks_multi, option.ShadowsocksUser{
				Name:     newUsers[i].Username,
				Password: newUsers[i].Name,
			})
			EditShadowsocksMultiUsers(c, []option.ShadowsocksUser{{
				Name:     newUsers[i].Username,
				Password: newUsers[i].Name,
			}}, delete)
		}

		if newUsers[i].AddToAll || newUsers[i].Shadowsocks_relay {
			shadowsocks_relay = append(shadowsocks_relay, option.ShadowsocksDestination{
				Name:     newUsers[i].Username,
				Password: newUsers[i].Name,
				ServerOptions: option.ServerOptions{
					Server:     newUsers[i].ServerAddress,
					ServerPort: newUsers[i].ServerPort,
				},
			})
			EditShadowsocksRelayUsers(c, []option.ShadowsocksDestination{{
				Name:     newUsers[i].Username,
				Password: newUsers[i].Name,
				ServerOptions: option.ServerOptions{
					Server:     newUsers[i].ServerAddress,
					ServerPort: newUsers[i].ServerPort,
				},
			}}, delete)
		}

		if newUsers[i].AddToAll || newUsers[i].Shadowtls {
			shadowtlsArr = append(shadowtlsArr, shadowtls.User{
				Name:     newUsers[i].Username,
				Password: newUsers[i].Name,
			})
			EditShadowtlsUsers(c, []shadowtls.User{{
				Name:     newUsers[i].Username,
				Password: newUsers[i].Name,
			}}, delete)
		}

		if newUsers[i].AddToAll || newUsers[i].Trojan {
			trojan = append(trojan, option.TrojanUser{
				Name:     newUsers[i].Username,
				Password: newUsers[i].Name,
			})
			EditTrojanUsers(c, []option.TrojanUser{{
				Name:     newUsers[i].Username,
				Password: newUsers[i].Name,
			}}, delete)
		}

		if newUsers[i].AddToAll || newUsers[i].Tuic {
			tuic = append(tuic, option.TUICUser{
				Name:     newUsers[i].Username,
				Password: newUsers[i].Name,
				UUID:     newUsers[i].Uuid,
			})
			EditTuicUsers(c, []option.TUICUser{{
				Name:     newUsers[i].Username,
				Password: newUsers[i].Name,
				UUID:     newUsers[i].Uuid,
			}}, delete)
		}

		if newUsers[i].AddToAll || newUsers[i].Vless {
			vless = append(vless, option.VLESSUser{
				Name: newUsers[i].Username,
				Flow: newUsers[i].Name,
				UUID: newUsers[i].Uuid,
			})
			EditVlessUsers(c, []option.VLESSUser{{
				Name: newUsers[i].Username,
				Flow: newUsers[i].Name,
				UUID: newUsers[i].Uuid,
			}}, delete)
		}

		if newUsers[i].AddToAll || newUsers[i].Vmess {
			vmess = append(vmess, option.VMessUser{
				Name: newUsers[i].Username,
				UUID: newUsers[i].Uuid,
			})
			EditVmessUsers(c, []option.VMessUser{{
				Name: newUsers[i].Username,
				UUID: newUsers[i].Uuid,
			}}, delete)
		}
	}

	//AddHysteria(c, hysteria)
	//AddHysteria2(c, hysteria2)
	//AddNaive(c, naive)
	//AddShadowsocksMulti(c, shadowsocks_multi)
	//AddShadowsocksRelay(c, shadowsocks_relay)
	//AddShadowTls(c, shadowtlsArr)
	//AddTrojan(c, trojan)
	//AddTuic(c, tuic)
	//AddVless(c, vless)
	//AddVmess(c, vmess)
}
