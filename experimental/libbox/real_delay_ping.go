package libbox

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"time"

	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/adapter"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/bufio"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
	"github.com/sagernet/sing/service"
	"golang.org/x/net/context"
)

func GetRealDelayPing(url string, config string, platformInterface PlatformInterface) int64 {
	C.ENCRYPTED_CONFIG = true
	return fetchDomesticPlatformInterface(url, config, platformInterface)
}

type OptionsEntry struct {
	content []byte
	path    string
	options option.Options
}

func readConfigAt(path string) (*OptionsEntry, error) {
	var (
		configContent []byte
		err           error
	)
	if path == "stdin" {
		configContent, err = io.ReadAll(os.Stdin)
	} else {
		configContent, err = os.ReadFile(path)
	}

	if err != nil {
		if C.ENCRYPTED_CONFIG {
			configContent = []byte(box.Decrypt(path))
			err = nil
		}
	} else {
		if C.ENCRYPTED_CONFIG {
			configContent = []byte(box.Decrypt(string(configContent)))
			err = nil
		}
	}

	if err != nil {
		return nil, E.Cause(err, "read config at ", path)
	}
	options, err := json.UnmarshalExtended[option.Options](configContent)
	if err != nil {
		return nil, E.Cause(err, "decode config at ", path)
	}
	return &OptionsEntry{
		content: configContent,
		path:    path,
		options: options,
	}, nil
}

func ReadEncryptedConfig(config string) ([]*OptionsEntry, error) {
	var optionsList []*OptionsEntry
	optionsEntry, err := readConfigAt(config)
	if err != nil {
		return nil, err
	}
	optionsList = append(optionsList, optionsEntry)
	sort.Slice(optionsList, func(i, j int) bool {
		return optionsList[i].path < optionsList[j].path
	})
	return optionsList, nil
}

func createDialer(instance *box.Box, network string, outboundTag string) (N.Dialer, error) {
	if outboundTag == "" {
		return instance.Outbound().Default(), nil
	} else {
		outbound, loaded := instance.Outbound().Outbound(outboundTag)
		if !loaded {
			return nil, E.New("outbound not found: ", outboundTag)
		}
		return outbound, nil
	}
}

func fetchDomesticPlatformInterface(url string, args string, platformInterface PlatformInterface) int64 {
	ctx := BaseContext(platformInterface)
	options, err := parseConfig(ctx, args)
	if err != nil {
		log.Error("RealDelay:-1")
		log.Error(err.Error())
		return -1
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ctx = service.ContextWith[adapter.PlatformInterface](ctx, &platformInterfaceWrapper{iif: platformInterface})

	instance, err := box.New(box.Options{
		Context: ctx,
		Options: options,
	})
	if err != nil {
		log.Error("RealDelay:-1")
		log.Error(err.Error())
		return -1
	}
	defer instance.Close()
	return fetchDomestic(url, instance)
}
func fetchDomestic(urll string, instance *box.Box) int64 {
	if instance == nil {
		return -1
	}

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSHandshakeTimeout: 5 * time.Second,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				dialer, err := createDialer(instance, network, "")
				if err != nil {
					log.Error(err.Error())
					return nil, err
				}
				return dialer.DialContext(ctx, network, M.ParseSocksaddr(addr))
			},
			ForceAttemptHTTP2: true,
		},
	}
	defer httpClient.CloseIdleConnections()

	parsedURL, err := url.Parse(urll)
	if err != nil {
		log.Error(err.Error())
		log.Error("RealDelay:-1")
		return -1
	}

	switch parsedURL.Scheme {
	case "":
		parsedURL.Scheme = "http"
		fallthrough
	case "http", "https":
		return fetchHTTPWithClient(httpClient, parsedURL)
	}
	return -1
}

func fetchHTTPWithClient(httpClient *http.Client, parsedURL *url.URL) int64 {
	request, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		log.Error(err.Error())
		return -1
	}
	request.Header.Add("User-Agent", "curl/7.88.0")
	start := time.Now()
	response, err := httpClient.Do(request)

	if response != nil {
		defer response.Body.Close()
		_, err = bufio.Copy(os.Stdout, response.Body)
		if errors.Is(err, io.EOF) {
			log.Error(err.Error())
			return -1
		}
		if err != nil {
			log.Error(err.Error())
			log.Error("RealDelay:-1")
			return -1
		}

		if response.StatusCode != http.StatusNoContent {
			log.Error("RealDelay:-1")
		}
		pingTime := time.Since(start).Milliseconds()
		log.Info("RealDelay:" + strconv.FormatInt(pingTime, 10))
		return pingTime
	} else {
		log.Error(err.Error())
		log.Error("RealDelay:-1")
		return -1
	}

}
