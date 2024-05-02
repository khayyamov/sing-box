package libbox

import (
	"errors"
	"fmt"
	box "github.com/sagernet/sing-box"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/bufio"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/json/badjson"
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

func GetRealDelayPing(config string, platformInterface PlatformInterface) int64 {
	C.ENCRYPTED_CONFIG = true
	return fetchDomesticPlatformInterface(config, platformInterface)
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

func readEncryptedConfigAndMerge(config string) (option.Options, error) {
	optionsList, err := ReadEncryptedConfig(config)
	if err != nil {
		return option.Options{}, err
	}
	if len(optionsList) == 1 {
		return optionsList[0].options, nil
	}
	var mergedMessage json.RawMessage
	for _, options := range optionsList {
		mergedMessage, err = badjson.MergeJSON(options.options.RawMessage, mergedMessage)
		if err != nil {
			return option.Options{}, E.Cause(err, "merge config at ", options.path)
		}
	}
	var mergedOptions option.Options
	err = mergedOptions.UnmarshalJSON(mergedMessage)
	if err != nil {
		return option.Options{}, E.Cause(err, "unmarshal merged config")
	}
	return mergedOptions, nil
}

func createPreStartedClientForApi(config string) (*box.Box, error) {
	options, err := readEncryptedConfigAndMerge(config)
	if err != nil {
		return nil, err
	}
	instance, err := box.New(box.Options{Options: options})
	if err != nil {
		return nil, E.Cause(err, "create service")
	}
	err = instance.PreStart()
	if err != nil {
		return nil, E.Cause(err, "start service")
	}
	return instance, nil
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

func fetchDomesticPlatformInterface(args string, platformInterface PlatformInterface) int64 {

	instance, errr := NewService(args, platformInterface)

	if errr != nil {
		log.Error("RealDelay:-1")
		log.Error(errr.Error())
		return -1
	}
	defer instance.Close()
	return fetchDomestic(instance)
}
func fetchDomestic(instance *BoxService) int64 {
	if instance != nil {
		if instance.instance != nil {
			log.Info("kilo 1")
			fmt.Println("kilo 1")
			httpClientt = &http.Client{
				Timeout: 5 * time.Second,
				Transport: &http.Transport{
					TLSHandshakeTimeout: 5 * time.Second,
					DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
						log.Info("kilo 2")
						fmt.Println("kilo 2")
						dialer, err := createDialer(instance.instance, network, "")
						if err != nil {
							log.Error(err.Error())
							return nil, err
						}
						log.Info("kilo 3")
						fmt.Println("kilo 3")
						return dialer.DialContext(ctx, network, M.ParseSocksaddr(addr))
					},
					ForceAttemptHTTP2: true,
				},
			}
			log.Info("kilo 4")
			fmt.Println("kilo 4")
			defer httpClientt.CloseIdleConnections()
			parsedURL, err := url.Parse("https://www.google.com/generate_204")
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
	log.Info("kilo 5")
	fmt.Println("kilo 5")
	request, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		log.Error(err.Error())
		return -1
	}
	log.Info("kilo 6")
	fmt.Println("kilo 6")
	request.Header.Add("User-Agent", "curl/7.88.0")
	start := time.Now()
	response, err := httpClientt.Do(request)
	log.Info("kilo 7")
	fmt.Println("kilo 7")
	defer response.Body.Close()
	_, err = bufio.Copy(os.Stdout, response.Body)
	log.Info("kilo 8")
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
}
