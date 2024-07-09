package inbound

import (
	"context"
	"github.com/sagernet/sing-box/adapter"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/experimental/libbox/platform"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
)

var TunPtr = option.Listable[*Tun]{}
var RedirectPtr = option.Listable[*Redirect]{}
var TProxyPtr = option.Listable[*TProxy]{}
var DirectPtr = option.Listable[*Direct]{}
var HTTPPtr = option.Listable[*HTTP]{}
var SocksPtr = option.Listable[*Socks]{}
var MixedPtr = option.Listable[*Mixed]{}

var TrojanPtr = option.Listable[*Trojan]{}
var VLESSPtr = option.Listable[*VLESS]{}
var VMessPtr = option.Listable[*VMess]{}

var NaivePtr = option.Listable[*Naive]{}
var ShadowTlsPtr = option.Listable[*ShadowTLS]{}
var ShadowsocksPtr = option.Listable[*Shadowsocks]{}
var ShadowsocksMultiPtr = option.Listable[*ShadowsocksMulti]{}
var ShadowsocksRelayPtr = option.Listable[*ShadowsocksRelay]{}
var TUICPtr = option.Listable[*TUIC]{}
var HysteriaPtr = option.Listable[*Hysteria]{}
var Hysteria2Ptr = option.Listable[*Hysteria2]{}

func New(ctx context.Context, router adapter.Router, logger log.ContextLogger, options option.Inbound, tag string, platformInterface platform.Interface) (adapter.Inbound, error) {
	if options.Type == "" {
		return nil, E.New("missing inbound type")
	}
	switch options.Type {
	//TODO: remove if need it
	case C.TypeTun:
		tun, err := NewTun(ctx, router, logger, tag, options.TunOptions, platformInterface)
		TunPtr = append(TunPtr, tun)
		return tun, err
	//case C.TypeRedirect:
	//	redirect := NewRedirect(ctx, router, logger, tag, options.RedirectOptions)
	//	RedirectPtr = append(RedirectPtr, redirect)
	//	return redirect, nil
	//case C.TypeTProxy:
	//	tproxy := NewTProxy(ctx, router, logger, tag, options.TProxyOptions)
	//	TProxyPtr = append(TProxyPtr, tproxy)
	//	return tproxy, nil
	//case C.TypeDirect:
	//	DirectPtr = NewDirect(ctx, router, logger, tag, options.DirectOptions)
	//	RedirectPtr = append(RedirectPtr, tun)
	//	return DirectPtr, nil
	//case C.TypeSOCKS:
	//	SocksPtr = NewSocks(ctx, router, logger, tag, options.SocksOptions)
	//	RedirectPtr = append(RedirectPtr, tun)
	//	return SocksPtr, nil
	//case C.TypeHTTP:
	//	HTTPPtr, err = NewHTTP(ctx, router, logger, tag, options.HTTPOptions)
	//	RedirectPtr = append(RedirectPtr, tun)
	//	return HTTPPtr, err
	//case C.TypeMixed:
	//	MixedPtr = NewMixed(ctx, router, logger, tag, options.MixedOptions)
	//	RedirectPtr = append(RedirectPtr, tun)
	//	return MixedPtr, nil

	case C.TypeVMess:
		vmess, err := NewVMess(ctx, router, logger, tag, options.VMessOptions)
		VMessPtr = append(VMessPtr, vmess)
		return vmess, err
	case C.TypeTrojan:
		trojan, err := NewTrojan(ctx, router, logger, tag, options.TrojanOptions)
		TrojanPtr = append(TrojanPtr, trojan)
		return trojan, err
	case C.TypeNaive:
		naive, err := NewNaive(ctx, router, logger, tag, options.NaiveOptions)
		NaivePtr = append(NaivePtr, naive)
		return naive, err
	case C.TypeShadowTLS:
		shadowtls, err := NewShadowTLS(ctx, router, logger, tag, options.ShadowTLSOptions)
		ShadowTlsPtr = append(ShadowTlsPtr, shadowtls)
		return shadowtls, err
	case C.TypeVLESS:
		vless, err := NewVLESS(ctx, router, logger, tag, options.VLESSOptions)
		VLESSPtr = append(VLESSPtr, vless)
		return vless, err

	case C.TypeShadowsocks:
		if len(options.ShadowsocksOptions.Users) > 0 && len(options.ShadowsocksOptions.Destinations) > 0 {
			return nil, E.New("users and destinations options must not be combined")
		}
		if len(options.ShadowsocksOptions.Users) > 0 {
			shadowsocks, err := newShadowsocksMulti(ctx, router, logger, tag, options.ShadowsocksOptions)
			ShadowsocksMultiPtr = append(ShadowsocksMultiPtr, shadowsocks)
			return shadowsocks, err
		} else if len(options.ShadowsocksOptions.Destinations) > 0 {
			shadowsocks, err := newShadowsocksRelay(ctx, router, logger, tag, options.ShadowsocksOptions)
			ShadowsocksRelayPtr = append(ShadowsocksRelayPtr, shadowsocks)
			return shadowsocks, err
		} else {
			shadowsocks, err := newShadowsocks(ctx, router, logger, tag, options.ShadowsocksOptions)
			ShadowsocksPtr = append(ShadowsocksPtr, shadowsocks)
			return shadowsocks, err
		}
	case C.TypeHysteria:
		hysteria, err := NewHysteria(ctx, router, logger, tag, options.HysteriaOptions)
		HysteriaPtr = append(HysteriaPtr, hysteria)
		return hysteria, err
	case C.TypeTUIC:
		tuic, err := NewTUIC(ctx, router, logger, tag, options.TUICOptions)
		TUICPtr = append(TUICPtr, tuic)
		return tuic, err
	case C.TypeHysteria2:
		hysteria2, err := NewHysteria2(ctx, router, logger, tag, options.Hysteria2Options)
		Hysteria2Ptr = append(Hysteria2Ptr, hysteria2)
		return hysteria2, err
	default:
		return nil, E.New("unknown inbound type: ", options.Type)
	}
}
