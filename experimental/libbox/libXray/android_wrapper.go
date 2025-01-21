//go:build android

package libXray

import (
	c "github.com/sagernet/sing-box/experimental/libbox/libXray/controller"
	"github.com/sagernet/sing-box/experimental/libbox/libXray/dns"
)

type DialerController interface {
	ProtectFd(int) bool
}

func InitDns(controller DialerController, server string) {
	dns.InitDns(server, func(fd uintptr) {
		controller.ProtectFd(int(fd))
	})
}

func ResetDns() {
	dns.ResetDns()
}

func RegisterDialerController(controller DialerController) {
	c.RegisterDialerController(func(fd uintptr) {
		controller.ProtectFd(int(fd))
	})
}

func RegisterListenerController(controller DialerController) {
	c.RegisterListenerController(func(fd uintptr) {
		controller.ProtectFd(int(fd))
	})
}
