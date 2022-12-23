package addrproviders

import (
	"context"
	"github.com/google/wire"
	"github.com/kelseyhightower/envconfig"
	"go-dyndns/util"
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
	config     *WebProviderConfig
	httpClient util.HttpClient
}

func newWebProvider(config *WebProviderConfig, httpClient util.HttpClient) *webProvider {
	return &webProvider{config: config, httpClient: httpClient}
}

func loadWebProviderConfig() (*WebProviderConfig, error) {
	var config WebProviderConfig
	if err := envconfig.Process("ipcheck", &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (w *webProvider) GetIPv4(ctx context.Context) (net.IP, error) {
	body, err := w.getBodyText(ctx, w.config.Url)
	if err != nil {
		return nil, err
	}

	startIdx := strings.Index(body, ":")
	if startIdx < 0 {
		return nil, InvalidResponseError
	}

	ipStr := strings.TrimSpace(body[startIdx+1:])
	if ip := net.ParseIP(ipStr); ip == nil {
		return nil, ParseIpError
	} else {
		return ip, nil
	}
}

func (w *webProvider) getBodyText(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	node, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	body := htmlquery.FindOne(node, "//body/text()")
	if body == nil {
		return "", InvalidResponseError
	}

	return body.Data, nil
}
