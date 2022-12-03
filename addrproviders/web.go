package addrproviders

import (
	"errors"
	"github.com/google/wire"
	"github.com/kelseyhightower/envconfig"
	"net"
	"net/http"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

var webSet = wire.NewSet(loadWebProviderConfig, newWebProvider)

type WebProviderConfig struct {
	Url string `required:"true"`
}

type webProvider struct {
	config *WebProviderConfig
}

func newWebProvider(config *WebProviderConfig) *webProvider {
	return &webProvider{config: config}
}

func loadWebProviderConfig() (*WebProviderConfig, error) {
	var config WebProviderConfig
	if err := envconfig.Process("web", &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (w *webProvider) GetIP() (net.IP, error) {
	body, err := getBodyText(w.config.Url)
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

func getBodyText(url string) (string, error) {
	resp, err := http.Get(url)
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
