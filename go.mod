module github.com/sagernet/sing-box

go 1.23

toolchain go1.23.4

require (
	github.com/caddyserver/certmagic v0.20.0
	github.com/cloudflare/circl v1.5.0
	github.com/cretz/bine v0.2.0
	github.com/go-chi/chi/v5 v5.1.0
	github.com/go-chi/render v1.0.3
	github.com/gofrs/uuid/v5 v5.3.0
	github.com/insomniacslk/dhcp v0.0.0-20231206064809-8c70d406f6d2
	github.com/libdns/alidns v1.0.3
	github.com/libdns/cloudflare v0.1.1
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/metacubex/tfo-go v0.0.0-20241006021335-daedaf0ca7aa
	github.com/mholt/acmez v1.2.0
	github.com/miekg/dns v1.1.62
	github.com/oschwald/maxminddb-golang v1.12.0
	github.com/sagernet/asc-go v0.0.0-20241217030726-d563060fe4e1
	github.com/sagernet/bbolt v0.0.0-20231014093535-ea5cb2fe9f0a
	github.com/sagernet/cloudflare-tls v0.0.0-20231208171750-a4483c1b7cd1
	github.com/sagernet/cors v1.2.1
	github.com/sagernet/fswatch v0.1.1
	github.com/sagernet/gomobile v0.1.4
	github.com/sagernet/gvisor v0.0.0-20241123041152-536d05261cff
	github.com/sagernet/quic-go v0.48.2-beta.1
	github.com/sagernet/reality v0.0.0-20230406110435-ee17307e7691
	github.com/sagernet/sing v0.6.0-beta.9
	github.com/sagernet/sing-dns v0.4.0-beta.1
	github.com/sagernet/sing-mux v0.3.0-alpha.1
	github.com/sagernet/sing-quic v0.4.0-beta.3
	github.com/sagernet/sing-shadowsocks v0.2.7
	github.com/sagernet/sing-shadowsocks2 v0.2.0
	github.com/sagernet/sing-shadowtls v0.2.0-alpha.2
	github.com/sagernet/sing-tun v0.6.0-beta.7
	github.com/sagernet/sing-vmess v0.2.0-beta.2
	github.com/sagernet/smux v0.0.0-20231208180855-7041f6ea79e7
	github.com/sagernet/utls v1.6.7
	github.com/sagernet/wireguard-go v0.0.1-beta.5
	github.com/sagernet/ws v0.0.0-20231204124109-acfe8907c854
	github.com/spf13/cobra v1.8.1
	github.com/stretchr/testify v1.10.0
	github.com/xtls/xray-core v1.8.25-0.20250103155023-de53a3b94e70
	go.uber.org/zap v1.27.0
	go4.org/netipx v0.0.0-20231129151722-fdeea329fbba
	golang.org/x/crypto v0.31.0
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56
	golang.org/x/mobile v0.0.0-20250103163954-5e5de4c85663
	golang.org/x/mod v0.22.0
	golang.org/x/net v0.33.0
	golang.org/x/sys v0.28.0
	golang.zx2c4.com/wireguard/wgctrl v0.0.0-20230429144221-925a1e7659e6
	google.golang.org/grpc v1.69.2
	google.golang.org/protobuf v1.36.1
	howett.net/plist v1.0.1
)

//replace github.com/sagernet/sing => ../sing

require (
	github.com/OmarTariq612/goech v0.0.0-20240405204721-8e2e1dafd3a0 // indirect
	github.com/ajg/form v1.5.1 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1 // indirect
	github.com/dgryski/go-metro v0.0.0-20211217172704-adc40b04c140 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/ghodss/yaml v1.0.1-0.20220118164431-d8423dcdf344 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/google/btree v1.1.3 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/pprof v0.0.0-20240528025155-186aa0362fba // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/hashicorp/yamux v0.1.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/josharian/native v1.1.0 // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/libdns/libdns v0.2.2 // indirect
	github.com/mdlayher/netlink v1.7.2 // indirect
	github.com/mdlayher/socket v0.4.1 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/onsi/ginkgo/v2 v2.19.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pierrec/lz4/v4 v4.1.14 // indirect
	github.com/pires/go-proxyproto v0.8.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/quic-go/qtls-go1-20 v0.4.1 // indirect
	github.com/refraction-networking/utls v1.6.7 // indirect
	github.com/riobard/go-bloom v0.0.0-20200614022211-cdc8013cb5b3 // indirect
	github.com/sagernet/netlink v0.0.0-20240612041022-b9a21c07ac6a // indirect
	github.com/sagernet/nftables v0.3.0-beta.4 // indirect
	github.com/seiflotfy/cuckoofilter v0.0.0-20240715131351-a2f2c23f1771 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/u-root/uio v0.0.0-20230220225925-ffce2a382923 // indirect
	github.com/v2fly/ss-bloomring v0.0.0-20210312155135-28617310f63e // indirect
	github.com/vishvananda/netlink v1.3.0 // indirect
	github.com/vishvananda/netns v0.0.4 // indirect
	github.com/xtls/quic-go v0.48.2 // indirect
	github.com/xtls/reality v0.0.0-20240712055506-48f0b2d5ed6d // indirect
	github.com/zeebo/blake3 v0.2.3 // indirect
	go.uber.org/mock v0.4.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/time v0.7.0 // indirect
	golang.org/x/tools v0.28.0 // indirect
	golang.zx2c4.com/wintun v0.0.0-20230126152724-0fa3db229ce2 // indirect
	golang.zx2c4.com/wireguard v0.0.0-20231211153847-12269c276173 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241015192408-796eee8c2d53 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gvisor.dev/gvisor v0.0.0-20240320123526-dc6abceb7ff0 // indirect
	lukechampine.com/blake3 v1.3.0 // indirect
)
