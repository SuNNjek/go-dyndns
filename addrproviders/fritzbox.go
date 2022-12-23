package addrproviders

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/kelseyhightower/envconfig"
	"go-dyndns/util"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/antchfx/xmlquery"
)

var fritzBoxSet = wire.NewSet(loadFritzBoxConfig, newFritzBoxProvider)

type fritzBoxConfig struct {
	Host string `default:"fritz.box"`
}

type fritzBoxProvider struct {
	config     *fritzBoxConfig
	httpClient util.HttpClient
}

func newFritzBoxProvider(config *fritzBoxConfig, httpClient util.HttpClient) *fritzBoxProvider {
	return &fritzBoxProvider{config: config, httpClient: httpClient}
}

func loadFritzBoxConfig() (*fritzBoxConfig, error) {
	var config fritzBoxConfig
	if err := envconfig.Process("fritzbox", &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (f *fritzBoxProvider) GetIPv4(ctx context.Context) (net.IP, error) {
	resp, err := f.makeExternalIpSoapRequest(ctx, f.config.Host)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ipStr, err := parseExternalIpSoapResponse(resp.Body)
	if err != nil {
		return nil, err
	}

	if ip := net.ParseIP(ipStr); ip == nil {
		return nil, ParseIpError
	} else {
		return ip, nil
	}
}

func (f *fritzBoxProvider) GetIPv6Prefix(ctx context.Context) (*util.IPv6Prefix, error) {
	resp, err := f.makeIpv6PrefixSoapRequest(ctx, f.config.Host)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	result, err := parseIpv6PrefixSoapResponse(resp.Body)
	if err != nil {
		return nil, err
	}

	if prefix, err := util.ParseIPv6Prefix(result); err != nil {
		return nil, err
	} else {
		return prefix, nil
	}
}

func (f *fritzBoxProvider) makeExternalIpSoapRequest(ctx context.Context, host string) (*http.Response, error) {
	url := fmt.Sprintf("http://%s:49000/igdupnp/control/WANIPConn1", host)
	body := `<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
	<s:Body>
		<u:GetExternalIPAddress xmlns:u="urn:schemas-upnp-org:service:WANIPConnection:1" />
	</s:Body>
</s:Envelope>`

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	defer request.Body.Close()

	request.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	request.Header.Add("SOAPAction", "urn:schemas-upnp-org:service:WANIPConnection:1#GetExternalIPAddress")

	return f.httpClient.Do(request)
}
func (f *fritzBoxProvider) makeIpv6PrefixSoapRequest(ctx context.Context, host string) (*http.Response, error) {
	url := fmt.Sprintf("http://%s:49000/igdupnp/control/WANIPConn1", host)
	body := `<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
	<s:Body>
		<u:X_AVM_DE_GetIPv6Prefix xmlns:u="urn:schemas-upnp-org:service:WANIPConnection:1" />
	</s:Body>
</s:Envelope>`

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	defer request.Body.Close()

	request.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	request.Header.Add("SOAPAction", "urn:schemas-upnp-org:service:WANIPConnection:1#X_AVM_DE_GetIPv6Prefix")

	return f.httpClient.Do(request)
}

func parseExternalIpSoapResponse(body io.Reader) (string, error) {
	root, err := xmlquery.Parse(body)
	if err != nil {
		return "", err
	}

	node := xmlquery.FindOne(root, "//NewExternalIPAddress/text()")
	if node == nil {
		return "", InvalidResponseError
	}

	return node.Data, nil
}

func parseIpv6PrefixSoapResponse(body io.Reader) (string, error) {
	root, err := xmlquery.Parse(body)
	if err != nil {
		return "", err
	}

	prefixNode := xmlquery.FindOne(root, "//NewIPv6Prefix/text()")
	if prefixNode == nil {
		return "", InvalidResponseError
	}

	prefixLengthNode := xmlquery.FindOne(root, "//NewPrefixLength/text()")
	if prefixLengthNode == nil {
		return "", InvalidResponseError
	}

	addrStr := fmt.Sprintf("%s/%s", prefixNode.Data, prefixLengthNode.Data)
	return addrStr, nil
}
