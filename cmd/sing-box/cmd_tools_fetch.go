package main

import (
	"context"
	"errors"
	box "github.com/sagernet/sing-box"
	C "github.com/sagernet/sing-box/constant"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing/common/bufio"
	M "github.com/sagernet/sing/common/metadata"

	"github.com/spf13/cobra"
)

var commandFetch = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch an URL",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := fetch(args)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	commandTools.AddCommand(commandFetch)
}

var httpClient *http.Client

func fetch(args []string) error {
	return fetchDomestic(args[0], false)
}
func GetRealPingApi(config string) error {
	return fetchDomestic(config, true)
}

func fetchDomestic(args string, runFromApi bool) error {
	instance, errr := &box.Box{}, errors.New("")
	if runFromApi {
		C.ENCRYPTED_CONFIG = true
		instance, errr = createPreStartedClientForApi(args)
	} else {
		instance, errr = createPreStartedClient()
	}
	if errr != nil {
		return errors.New("RelayPing:-1")
	}
	defer instance.Close()
	httpClient = &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				dialer, err := createDialer(instance, network, commandToolsFlagOutbound)
				if err != nil {
					return nil, err
				}
				return dialer.DialContext(ctx, network, M.ParseSocksaddr(addr))
			},
			ForceAttemptHTTP2: true,
		},
	}
	defer httpClient.CloseIdleConnections()
	parsedURL, err := url.Parse(args)
	if err != nil {
		return err
	}
	switch parsedURL.Scheme {
	case "":
		parsedURL.Scheme = "http"
		fallthrough
	case "http", "https":
		err = fetchHTTP(parsedURL)
		if err != nil {
			return err
		}
	}
	return nil
}

func fetchHTTP(parsedURL *url.URL) error {
	request, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return err
	}
	request.Header.Add("User-Agent", "curl/7.88.0")
	start := time.Now()
	response, err := httpClient.Do(request)
	if err != nil {
		log.Error("RealDelay:-1")
		return err
	} else {
		if response.StatusCode != http.StatusNoContent {
			log.Error("RealDelay:-1")
		}
		log.Info("RealDelay:" + strconv.FormatInt(time.Since(start).Milliseconds(), 10))
	}
	defer response.Body.Close()
	_, err = bufio.Copy(os.Stdout, response.Body)
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}
