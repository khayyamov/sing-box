package inbound

import (
	"context"
	"github.com/sagernet/sing-box/api/constant"
	"github.com/sagernet/sing-box/api/db"
	"net"
	"os"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/common/dialer"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-shadowtls"
	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/auth"
	N "github.com/sagernet/sing/common/network"
)

type ShadowTLS struct {
	myInboundAdapter
	Service *shadowtls.Service
}

func NewShadowTLS(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.ShadowTLSInboundOptions) (*ShadowTLS, error) {
	if !constant.DbEnable {
		if len(options.Users) > 0 {
			users, _ := db.ConvertProtocolModelToDbUser(options.Users)
			db.GetDb().EditInRamUsers(users, false)
		}
	} else {
		dbUsers, _ := db.GetDb().GetShadowtlsUsers()
		for i := range dbUsers {
			options.Users = append(options.Users, option.ShadowTLSUser{
				Name:     dbUsers[i].Name,
				Password: dbUsers[i].Password,
			})
		}
		if len(dbUsers) > 0 {
			users, _ := db.ConvertProtocolModelToDbUser(dbUsers)
			db.GetDb().EditInRamUsers(users, false)
		}
	}
	inbound := &ShadowTLS{
		myInboundAdapter: myInboundAdapter{
			protocol:      C.TypeShadowTLS,
			network:       []string{N.NetworkTCP},
			ctx:           ctx,
			router:        router,
			logger:        logger,
			tag:           tag,
			ListenOptions: options.ListenOptions,
		},
	}

	if options.Version == 0 {
		options.Version = 1
	}

	var handshakeForServerName map[string]shadowtls.HandshakeConfig
	if options.Version > 1 {
		handshakeForServerName = make(map[string]shadowtls.HandshakeConfig)
		for serverName, serverOptions := range options.HandshakeForServerName {
			handshakeDialer, err := dialer.New(router, serverOptions.DialerOptions)
			if err != nil {
				return nil, err
			}
			handshakeForServerName[serverName] = shadowtls.HandshakeConfig{
				Server: serverOptions.ServerOptions.Build(),
				Dialer: handshakeDialer,
			}
		}
	}
	handshakeDialer, err := dialer.New(router, options.Handshake.DialerOptions)
	if err != nil {
		return nil, err
	}
	service, err := shadowtls.NewService(shadowtls.ServiceConfig{
		Version:  options.Version,
		Password: options.Password,
		Users: common.Map(options.Users, func(it option.ShadowTLSUser) shadowtls.User {
			return (shadowtls.User)(it)
		}),
		Handshake: shadowtls.HandshakeConfig{
			Server: options.Handshake.ServerOptions.Build(),
			Dialer: handshakeDialer,
		},
		HandshakeForServerName: handshakeForServerName,
		StrictMode:             options.StrictMode,
		Handler:                adapter.NewUpstreamContextHandler(inbound.newConnection, nil, inbound),
		Logger:                 logger,
	})
	if err != nil {
		return nil, err
	}
	inbound.Service = service
	inbound.connHandler = inbound
	return inbound, nil
}

func (h *ShadowTLS) NewConnection(ctx context.Context, conn net.Conn, metadata adapter.InboundContext) error {
	return h.Service.NewConnection(adapter.WithContext(log.ContextWithNewID(ctx), &metadata), conn, adapter.UpstreamMetadata(metadata))
}

func (h *ShadowTLS) newConnection(ctx context.Context, conn net.Conn, metadata adapter.InboundContext) error {
	if userName, _ := auth.UserFromContext[string](ctx); userName != "" {
		metadata.User = userName
		h.logger.InfoContext(ctx, "[", userName, "] inbound connection to ", metadata.Destination)
	} else {
		h.logger.InfoContext(ctx, "inbound connection to ", metadata.Destination)
		return os.ErrInvalid
	}
	return h.router.RouteConnection(ctx, conn, metadata)
}
