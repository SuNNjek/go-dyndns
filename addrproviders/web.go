package addrproviders

import (
	"errors"
	"github.com/google/wire"
	"github.com/kelseyhightower/envconfig"
	"go-dyndns/util"
	"net"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

var webSet = wire.NewSet(loadWebProviderConfig, newWebProvider)

type WebProviderConfig struct {
	Url string `required:"true"`
}

type webProvider struct {
	config     *WebProviderConfig
	httpClient util.HttpClient
}

func newWebProvider(config *WebProviderConfig, httpClient util.HttpClient) *webProvider {
	return &webProvider{config: config, httpClient: httpClient}
}

func loadWebProviderConfig() (*WebProviderConfig, error) {
	var config WebProviderConfig
	if err := envconfig.Process("web", &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (w *webProvider) GetIP() (net.IP, error) {
	body, err := w.getBodyText(w.config.Url)
	if err != nil {
		return nil, err
	}

	startIdx := strings.Index(body, ":")
	if startIdx < 0 {
		return nil, errors.New("no IP returned in response")
	}

	ipStr := strings.TrimSpace(body[startIdx+1:])
	if ip := net.ParseIP(ipStr); ip == nil {
		return nil, errors.New("failed to parse IP")
	} else {
		return ip, nil
	}
}

func (w *webProvider) getBodyText(url string) (string, error) {
	resp, err := w.httpClient.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	node, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	body, err := htmlquery.Query(node, "//body/text()")
	if err != nil {
		return "", err
	}

	return body.Data, nil
}
