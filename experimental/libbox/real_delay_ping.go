package libbox

import (
	"errors"
	box "github.com/sagernet/sing-box"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/bufio"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
	"golang.org/x/net/context"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"time"
)

var httpClientt *http.Client

func GetRealDelayPing(url string, config string, platformInterface PlatformInterface) int64 {
	log.Error("Crash 0")
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
		return instance.Router().DefaultOutbound(N.NetworkName(network))
	} else {
		outbound, loaded := instance.Router().Outbound(outboundTag)
		if !loaded {
			return nil, E.New("outbound not found: ", outboundTag)
		}
		return outbound, nil
	}
}

func fetchDomesticPlatformInterface(url string, args string, platformInterface PlatformInterface) int64 {

	instance, errr := NewService(args, platformInterface)

	if errr != nil {
		log.Error("RealDelay:-1")
		log.Error(errr.Error())
		return -1
	}
	defer instance.Close()
	return fetchDomestic(url, instance)
}
func fetchDomestic(urll string, instance *BoxService) int64 {
	if instance != nil {
		if instance.instance != nil {
			httpClientt = &http.Client{
				Timeout: 5 * time.Second,
				Transport: &http.Transport{
					TLSHandshakeTimeout: 5 * time.Second,
					DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
						dialer, err := createDialer(instance.instance, network, "")
						if err != nil {
							log.Error(err.Error())
							return nil, err
						}
						return dialer.DialContext(ctx, network, M.ParseSocksaddr(addr))
					},
					ForceAttemptHTTP2: true,
				},
			}
			defer httpClientt.CloseIdleConnections()
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
				return fetchHTTP(parsedURL)
			}
			return -1
		} else {
			return -1
		}
	} else {
		return -1
	}
}

func fetchHTTP(parsedURL *url.URL) int64 {
	request, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		log.Error(err.Error())
		return -1
	}
	request.Header.Add("User-Agent", "curl/7.88.0")
	start := time.Now()
	response, err := httpClientt.Do(request)

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
		} else {
			if response.StatusCode != http.StatusNoContent {
				log.Error("RealDelay:-1")
			}
			pingTime := time.Since(start).Milliseconds()
			log.Info("RealDelay:" + strconv.FormatInt(pingTime, 10))
			return pingTime
		}
	} else {
		log.Error(err.Error())
		log.Error("RealDelay:-1")
		return -1
	}

}
