package inbound

import (
	"context"
	"github.com/sagernet/sing-box/api/constant"
	"github.com/sagernet/sing-box/api/db"
	"net"
	"os"
	"time"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/common/mux"
	"github.com/sagernet/sing-box/common/uot"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-shadowsocks/shadowaead_2022"
	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/auth"
	"github.com/sagernet/sing/common/buf"
	F "github.com/sagernet/sing/common/format"
	N "github.com/sagernet/sing/common/network"
)

var (
	_ adapter.Inbound           = (*ShadowsocksRelay)(nil)
	_ adapter.InjectableInbound = (*ShadowsocksRelay)(nil)
)

type ShadowsocksRelay struct {
	myInboundAdapter
	Service      *shadowaead_2022.RelayService[string]
	Destinations map[string]option.ShadowsocksDestination
}

func newShadowsocksRelay(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.ShadowsocksInboundOptions) (*ShadowsocksRelay, error) {
	if !constant.DbEnable {
		if len(options.Users) > 0 {
			users, _ := db.ConvertProtocolModelToDbUser(options.Users)
			db.GetDb().EditInRamUsers(users, false)
		}
	} else {
		dbUsers, _ := db.GetDb().GetShadowsocksRelayUsers()
		dbUsers = append(dbUsers, options.Destinations...)
		if len(dbUsers) > 0 {
			options.Destinations = dbUsers
			users, _ := db.ConvertProtocolModelToDbUser(dbUsers)
			db.GetDb().EditInRamUsers(users, false)
		}
	}
	inbound := &ShadowsocksRelay{
		myInboundAdapter: myInboundAdapter{
			protocol:      C.TypeShadowsocks,
			network:       options.Network.Build(),
			ctx:           ctx,
			router:        uot.NewRouter(router, logger),
			logger:        logger,
			tag:           tag,
			ListenOptions: options.ListenOptions,
		},
		Destinations: map[string]option.ShadowsocksDestination{},
	}
	for _, user := range options.Destinations {
		inbound.Destinations[user.Name] = user
	}
	inbound.connHandler = inbound
	inbound.packetHandler = inbound
	var err error
	inbound.router, err = mux.NewRouterWithOptions(inbound.router, logger, common.PtrValueOrDefault(options.Multiplex))
	if err != nil {
		return nil, err
	}
	var udpTimeout time.Duration
	if options.UDPTimeout != 0 {
		udpTimeout = time.Duration(options.UDPTimeout)
	} else {
		udpTimeout = C.UDPTimeout
	}
	service, err := shadowaead_2022.NewRelayServiceWithPassword[string](
		options.Method,
		options.Password,
		int64(udpTimeout.Seconds()),
		adapter.NewUpstreamContextHandler(inbound.newConnection, inbound.newPacketConnection, inbound),
	)
	if err != nil {
		return nil, err
	}
	err = service.UpdateUsersWithPasswords(common.MapIndexedString(options.Destinations, func(index any, user option.ShadowsocksDestination) string {
		return user.Name
	}), common.Map(options.Destinations, func(user option.ShadowsocksDestination) string {
		return user.Password
	}), common.Map(options.Destinations, option.ShadowsocksDestination.Build))
	if err != nil {
		return nil, err
	}
	inbound.Service = service
	inbound.packetUpstream = service
	return inbound, err
}

func (h *ShadowsocksRelay) NewConnection(ctx context.Context, conn net.Conn, metadata adapter.InboundContext) error {
	return h.Service.NewConnection(adapter.WithContext(log.ContextWithNewID(ctx), &metadata), conn, adapter.UpstreamMetadata(metadata))
}

func (h *ShadowsocksRelay) NewPacket(ctx context.Context, conn N.PacketConn, buffer *buf.Buffer, metadata adapter.InboundContext) error {
	return h.Service.NewPacket(adapter.WithContext(ctx, &metadata), conn, buffer, adapter.UpstreamMetadata(metadata))
}

func (h *ShadowsocksRelay) NewPacketConnection(ctx context.Context, conn N.PacketConn, metadata adapter.InboundContext) error {
	return os.ErrInvalid
}

func (h *ShadowsocksRelay) newConnection(ctx context.Context, conn net.Conn, metadata adapter.InboundContext) error {
	destinationIndex, loaded := auth.UserFromContext[string](ctx)
	if !loaded {
		return os.ErrInvalid
	}
	destination := h.Destinations[destinationIndex].Name
	if destination == "" {
		destination = F.ToString(destinationIndex)
	} else {
		metadata.User = destination
		return os.ErrInvalid
	}
	h.logger.InfoContext(ctx, "[", destination, "] inbound connection to ", metadata.Destination)
	return h.router.RouteConnection(ctx, conn, metadata)
}

func (h *ShadowsocksRelay) newPacketConnection(ctx context.Context, conn N.PacketConn, metadata adapter.InboundContext) error {
	destinationIndex, loaded := auth.UserFromContext[string](ctx)
	if !loaded {
		return os.ErrInvalid
	}
	destination := h.Destinations[destinationIndex].Name
	if destination == "" {
		destination = F.ToString(destinationIndex)
	} else {
		metadata.User = destination
		return os.ErrInvalid
	}
	ctx = log.ContextWithNewID(ctx)
	h.logger.InfoContext(ctx, "[", destination, "] inbound packet connection from ", metadata.Source)
	h.logger.InfoContext(ctx, "[", destination, "] inbound packet connection to ", metadata.Destination)
	return h.router.RoutePacketConnection(ctx, conn, metadata)
}
