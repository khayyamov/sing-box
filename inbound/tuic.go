package inbound

import (
	"context"
	"github.com/gofrs/uuid/v5"
	"github.com/sagernet/sing-box/api/constant"
	"github.com/sagernet/sing-box/api/db"
	"net"
	"time"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/common/tls"
	"github.com/sagernet/sing-box/common/uot"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-quic/tuic"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
	N "github.com/sagernet/sing/common/network"
)

var _ adapter.Inbound = (*TUIC)(nil)

type TUIC struct {
	myInboundAdapter
	tlsConfig    tls.ServerConfig
	Service      *tuic.Service[int]
	Users        []option.TUICUser
	UserNameList []string
}

func NewTUIC(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.TUICInboundOptions) (*TUIC, error) {
	if !constant.DbEnable {
		if len(options.Users) > 0 {
			users, _ := db.ConvertProtocolModelToDbUser(options.Users)
			db.GetDb().EditInRamUsers(users, false)
		}
	} else {
		dbUsers, _ := db.GetDb().GetTuicUsers()
		dbUsers = append(dbUsers, options.Users...)
		if len(dbUsers) > 0 {
			options.Users = dbUsers
			users, _ := db.ConvertProtocolModelToDbUser(dbUsers)
			db.GetDb().EditInRamUsers(users, false)
		}
	}
	options.UDPFragmentDefault = true
	if options.TLS == nil || !options.TLS.Enabled {
		return nil, C.ErrTLSRequired
	}
	tlsConfig, err := tls.NewServer(ctx, logger, common.PtrValueOrDefault(options.TLS))
	if err != nil {
		return nil, err
	}
	inbound := &TUIC{
		myInboundAdapter: myInboundAdapter{
			protocol:      C.TypeTUIC,
			network:       []string{N.NetworkUDP},
			ctx:           ctx,
			router:        uot.NewRouter(router, logger),
			logger:        logger,
			tag:           tag,
			ListenOptions: options.ListenOptions,
		},
		tlsConfig: tlsConfig,
		Users:     options.Users,
	}
	var udpTimeout time.Duration
	if options.UDPTimeout != 0 {
		udpTimeout = time.Duration(options.UDPTimeout)
	} else {
		udpTimeout = C.UDPTimeout
	}
	service, err := tuic.NewService[int](tuic.ServiceOptions{
		Context:           ctx,
		Logger:            logger,
		TLSConfig:         tlsConfig,
		CongestionControl: options.CongestionControl,
		AuthTimeout:       time.Duration(options.AuthTimeout),
		ZeroRTTHandshake:  options.ZeroRTTHandshake,
		Heartbeat:         time.Duration(options.Heartbeat),
		UDPTimeout:        udpTimeout,
		Handler:           adapter.NewUpstreamHandler(adapter.InboundContext{}, inbound.newConnection, inbound.newPacketConnection, nil),
	})
	if err != nil {
		return nil, err
	}
	var userList []int
	var userNameList []string
	var userUUIDList [][16]byte
	var userPasswordList []string
	for index, user := range options.Users {
		if user.UUID == "" {
			return nil, E.New("missing uuid for user ", user.UUID)
		}
		userUUID, err := uuid.FromString(user.UUID)
		if err != nil {
			return nil, E.Cause(err, "invalid uuid for user ", user.UUID)
		}
		userList = append(userList, index)
		userNameList = append(userNameList, user.Name)
		userUUIDList = append(userUUIDList, userUUID)
		userPasswordList = append(userPasswordList, user.Password)
	}
	service.UpdateUsers(userList, userUUIDList, userPasswordList)
	inbound.Service = service
	inbound.UserNameList = userNameList
	return inbound, nil
}

func (h *TUIC) newConnection(ctx context.Context, conn net.Conn, metadata adapter.InboundContext) error {
	h.logger.InfoContext(ctx, "inbound connection from ", metadata.Source)
	h.logger.InfoContext(ctx, "inbound connection to ", metadata.Destination)
	return h.router.RouteConnection(ctx, conn, metadata)
}

func (h *TUIC) newPacketConnection(ctx context.Context, conn N.PacketConn, metadata adapter.InboundContext) error {
	h.logger.InfoContext(ctx, "inbound packet connection from ", metadata.Source)
	h.logger.InfoContext(ctx, "inbound packet connection to ", metadata.Destination)
	return h.router.RoutePacketConnection(ctx, conn, metadata)
}

func (h *TUIC) Start() error {
	if h.tlsConfig != nil {
		err := h.tlsConfig.Start()
		if err != nil {
			return err
		}
	}
	packetConn, err := h.myInboundAdapter.ListenUDP()
	if err != nil {
		return err
	}
	return h.Service.Start(packetConn)
}

func (h *TUIC) Close() error {
	return common.Close(
		&h.myInboundAdapter,
		h.tlsConfig,
		common.PtrOrNil(h.Service),
	)
}
