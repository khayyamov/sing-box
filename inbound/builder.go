package inbound

import (
	"context"
	"errors"
	"github.com/sagernet/sing-box/adapter"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/experimental/libbox/platform"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
)

var TunPtr = &Tun{}
var RedirectPtr = &Redirect{}
var TProxyPtr = &TProxy{}
var DirectPtr = &Direct{}
var HTTPPtr = &HTTP{}
var SocksPtr = &Socks{}
var MixedPtr = &Mixed{}
var TrojanPtr = &Trojan{}
var NaivePtr = &Naive{}
var ShadowTlsPtr = &ShadowTLS{}
var VLESSPtr = &VLESS{}
var VMessPtr = &VMess{}

var ShadowsocksPtr = &Shadowsocks{}
var TUICPtr = &TUIC{}
var HysteriaPtr = &Hysteria{}
var Hysteria2Ptr = &Hysteria2{}

// TODO: get users from db for each protocol in first run
func New(ctx context.Context, router adapter.Router, logger log.ContextLogger, options option.Inbound, platformInterface platform.Interface) (adapter.Inbound, error) {
	if options.Type == "" {
		return nil, E.New("missing inbound type")
	}
	var err = errors.New("")
	switch options.Type {
	case C.TypeTun:
		TunPtr, err = NewTun(ctx, router, logger, options.Tag, options.TunOptions, platformInterface)
		return TunPtr, err
	case C.TypeRedirect:
		RedirectPtr = NewRedirect(ctx, router, logger, options.Tag, options.RedirectOptions)
		return RedirectPtr, nil
	case C.TypeTProxy:
		TProxyPtr = NewTProxy(ctx, router, logger, options.Tag, options.TProxyOptions)
		return TProxyPtr, nil
	case C.TypeDirect:
		DirectPtr = NewDirect(ctx, router, logger, options.Tag, options.DirectOptions)
		return DirectPtr, nil
	case C.TypeSOCKS:
		SocksPtr = NewSocks(ctx, router, logger, options.Tag, options.SocksOptions)
		return SocksPtr, nil
	case C.TypeHTTP:
		HTTPPtr, err = NewHTTP(ctx, router, logger, options.Tag, options.HTTPOptions)
		return HTTPPtr, err
	case C.TypeMixed:
		MixedPtr = NewMixed(ctx, router, logger, options.Tag, options.MixedOptions)
		return MixedPtr, nil
	case C.TypeVMess:
		VMessPtr, err = NewVMess(ctx, router, logger, options.Tag, options.VMessOptions)
		return VMessPtr, err
	case C.TypeTrojan:
		TrojanPtr, err = NewTrojan(ctx, router, logger, options.Tag, options.TrojanOptions)
		return TrojanPtr, err
	case C.TypeNaive:
		NaivePtr, err = NewNaive(ctx, router, logger, options.Tag, options.NaiveOptions)
		return NaivePtr, err
	case C.TypeShadowTLS:
		ShadowTlsPtr, err = NewShadowTLS(ctx, router, logger, options.Tag, options.ShadowTLSOptions)
		return ShadowTlsPtr, err
	case C.TypeVLESS:
		VLESSPtr, err = NewVLESS(ctx, router, logger, options.Tag, options.VLESSOptions)
		return VLESSPtr, err

	case C.TypeShadowsocks:
		ShadowsocksPtr, err = newShadowsocks(ctx, router, logger, options.Tag, options.ShadowsocksOptions)
		return ShadowsocksPtr, err
	case C.TypeHysteria:
		HysteriaPtr, err = NewHysteria(ctx, router, logger, options.Tag, options.HysteriaOptions)
		return HysteriaPtr, err
	case C.TypeTUIC:
		TUICPtr, err = NewTUIC(ctx, router, logger, options.Tag, options.TUICOptions)
		return TUICPtr, err
	case C.TypeHysteria2:
		Hysteria2Ptr, err = NewHysteria2(ctx, router, logger, options.Tag, options.Hysteria2Options)
		return Hysteria2Ptr, err
	default:
		return nil, E.New("unknown inbound type: ", options.Type)
	}
}
