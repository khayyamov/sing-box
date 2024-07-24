package rest

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	global_import "github.com/sagernet/sing-box/api/rest/rq"
	"github.com/sagernet/sing-box/api/utils"
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
	var rqArr []global_import.GlobalModel

	bodyAsByteArray, _ := io.ReadAll(c.Request.Body)
	jsonBody := string(bodyAsByteArray)

	err := json.Unmarshal([]byte(jsonBody), &rqArr)
	if err == nil {
		if len(rqArr) > 0 {
			var prettyJSON bytes.Buffer
			error := json.Indent(&prettyJSON, []byte(jsonBody), "", "\t")
			if error != nil {
				utils.ApiLogError("http.StatusBadRequest: " + error.Error())
			}
			utils.ApiLogInfo("Api Request Body: " + string(prettyJSON.Bytes()))
			EditGroupUsers(c, rqArr, delete)
		}
	} else {
		utils.ApiLogError("http.StatusBadRequest: " + err.Error())
		c.JSON(http.StatusBadRequest, err)
	}
}

func EditGroupUsers(c *gin.Context, newUsers []global_import.GlobalModel, delete bool) {
	EditVlessUsers(c, newUsers, delete)
	EditVmessUsers(c, newUsers, delete)
	EditTrojanUsers(c, newUsers, delete)
	EditTuicUsers(c, newUsers, delete)
	EditShadowtlsUsers(c, newUsers, delete)
	EditNaiveUsers(c, newUsers, delete)
	EditShadowsocksRelayUsers(c, newUsers, delete)
	EditShadowsocksMultiUsers(c, newUsers, delete)
	EditHysteriaUsers(c, newUsers, delete)
	EditHysteria2Users(c, newUsers, delete)
}

//func EditGroupUsers(c *gin.Context, newUsers []global_import.GlobalModel, delete bool) {
//	//hysteria := make([]option.HysteriaUser, len(newUsers))
//	//hysteria2 := make([]option.Hysteria2User, len(newUsers))
//	//naive := make([]auth.User, len(newUsers))
//	//shadowsocks_multi := make([]option.ShadowsocksUser, len(newUsers))
//	//shadowsocks_relay := make([]option.ShadowsocksDestination, len(newUsers))
//	//shadowtlsArr := make([]shadowtls.User, len(newUsers))
//	//trojan := make([]option.TrojanUser, len(newUsers))
//	//tuic := make([]option.TUICUser, len(newUsers))
//	//vless := make([]option.VLESSUser, len(newUsers))
//	//vmess := make([]option.VMessUser, len(newUsers))
//	for i := range newUsers {
//		if newUsers[i].AddToAll || newUsers[i].Hysteria {
//			//hysteria = append(hysteria, option.HysteriaUser{
//			//	UUID:       newUsers[i].UUID,
//			//	Auth:       []byte(newUsers[i].Auth),
//			//	AuthString: newUsers[i].AuthString,
//			//})
//			EditHysteriaUsers(c, []option.HysteriaUser{{
//				Name:       newUsers[i].UUID,
//				Auth:       []byte(newUsers[i].Auth),
//				AuthString: newUsers[i].AuthString,
//			}}, delete)
//		}
//
//		if newUsers[i].AddToAll || newUsers[i].Hysteria2 {
//			//hysteria2 = append(hysteria2, option.Hysteria2User{
//			//	UUID:     newUsers[i].UUID,
//			//	Password: newUsers[i].UUID,
//			//})
//			EditHysteria2Users(c, []option.Hysteria2User{{
//				Name:     newUsers[i].UUID,
//				Password: newUsers[i].Password,
//			}}, delete)
//		}
//
//		if newUsers[i].AddToAll || newUsers[i].Naive {
//			//naive = append(naive, auth.User{
//			//	UUID: newUsers[i].UUID,
//			//	Password: newUsers[i].UUID,
//			//})
//			EditNaiveUsers(c, []auth.User{{
//				Username: newUsers[i].UUID,
//				Password: newUsers[i].Password,
//			}}, delete)
//		}
//
//		if newUsers[i].AddToAll || newUsers[i].Shadowsocks_multi {
//			//shadowsocks_multi = append(shadowsocks_multi, option.ShadowsocksUser{
//			//	UUID:     newUsers[i].UUID,
//			//	Password: newUsers[i].UUID,
//			//})
//			EditShadowsocksMultiUsers(c, []option.ShadowsocksUser{{
//				Name:     newUsers[i].UUID,
//				Password: newUsers[i].Password,
//			}}, delete)
//		}
//
//		if newUsers[i].AddToAll || newUsers[i].Shadowsocks_relay {
//			//shadowsocks_relay = append(shadowsocks_relay, option.ShadowsocksDestination{
//			//	UUID:     newUsers[i].UUID,
//			//	Password: newUsers[i].UUID,
//			//	ServerOptions: option.ServerOptions{
//			//		Server:     newUsers[i].ServerAddress,
//			//		ServerPort: newUsers[i].ServerPort,
//			//	},
//			//})
//			EditShadowsocksRelayUsers(c, []option.ShadowsocksDestination{{
//				Name:     newUsers[i].UUID,
//				Password: newUsers[i].Password,
//				ServerOptions: option.ServerOptions{
//					Server:     newUsers[i].ServerAddress,
//					ServerPort: newUsers[i].ServerPort,
//				},
//			}}, delete)
//		}
//
//		if newUsers[i].AddToAll || newUsers[i].Shadowtls {
//			//shadowtlsArr = append(shadowtlsArr, shadowtls.User{
//			//	UUID:     newUsers[i].UUID,
//			//	Password: newUsers[i].UUID,
//			//})
//			EditShadowtlsUsers(c, []shadowtls.User{{
//				Name:     newUsers[i].UUID,
//				Password: newUsers[i].Password,
//			}}, delete)
//		}
//
//		if newUsers[i].AddToAll || newUsers[i].Trojan {
//			//trojan = append(trojan, option.TrojanUser{
//			//	UUID:     newUsers[i].UUID,
//			//	Password: newUsers[i].UUID,
//			//})
//			EditTrojanUsers(c, []option.TrojanUser{{
//				Name:     newUsers[i].UUID,
//				Password: newUsers[i].Password,
//			}}, delete)
//		}
//
//		if newUsers[i].AddToAll || newUsers[i].Tuic {
//			//tuic = append(tuic, option.TUICUser{
//			//	UUID:     newUsers[i].UUID,
//			//	Password: newUsers[i].UUID,
//			//	UUID:     newUsers[i].Uuid,
//			//})
//			EditTuicUsers(c, []option.TUICUser{{
//				Name:     newUsers[i].UUID,
//				Password: newUsers[i].Password,
//				UUID:     newUsers[i].UUID,
//			}}, delete)
//		}
//
//		if newUsers[i].AddToAll || newUsers[i].Vless {
//			//vless = append(vless, option.VLESSUser{
//			//	UUID: newUsers[i].UUID,
//			//	Flow: newUsers[i].UUID,
//			//	UUID: newUsers[i].Uuid,
//			//})
//			EditVlessUsers(c, []option.VLESSUser{{
//				Name: newUsers[i].UUID,
//				Flow: newUsers[i].Flow,
//				UUID: newUsers[i].UUID,
//			}}, delete)
//		}
//
//		if newUsers[i].AddToAll || newUsers[i].Vmess {
//			//vmess = append(vmess, option.VMessUser{
//			//	UUID: newUsers[i].UUID,
//			//	UUID: newUsers[i].Uuid,
//			//})
//			EditVmessUsers(c, []option.VMessUser{{
//				Name: newUsers[i].UUID,
//				UUID: newUsers[i].UUID,
//			}}, delete)
//		}
//	}
//
//	//AddHysteria(c, hysteria)
//	//AddHysteria2(c, hysteria2)
//	//AddNaive(c, naive)
//	//AddShadowsocksMulti(c, shadowsocks_multi)
//	//AddShadowsocksRelay(c, shadowsocks_relay)
//	//AddShadowTls(c, shadowtlsArr)
//	//AddTrojan(c, trojan)
//	//AddTuic(c, tuic)
//	//AddVless(c, vless)
//	//AddVmess(c, vmess)
//}
