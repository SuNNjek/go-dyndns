package addrproviders

import (
	"errors"
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

func (f *fritzBoxProvider) GetIP() (net.IP, error) {
	resp, err := f.makeExternalIpSoapRequest(f.config.Host)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ipStr, err := parseExternalIpSoapResponse(resp.Body)
	if err != nil {
		return nil, err
	}

	if ip := net.ParseIP(ipStr); ip == nil {
		return nil, errors.New("failed to parse IP")
	} else {
		return ip, nil
	}
}

func (f *fritzBoxProvider) makeExternalIpSoapRequest(host string) (*http.Response, error) {
	url := fmt.Sprintf("http://%s:49000/igdupnp/control/WANIPConn1", host)
	body := `<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
	<s:Body>
		<u:GetExternalIPAddress xmlns:u="urn:schemas-upnp-org:service:WANIPConnection:1" />
	</s:Body>
</s:Envelope>`

	request, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	defer request.Body.Close()

	request.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	request.Header.Add("SOAPAction", "urn:schemas-upnp-org:service:WANIPConnection:1#GetExternalIPAddress")

	return f.httpClient.Do(request)
}

func parseExternalIpSoapResponse(body io.Reader) (string, error) {
	root, err := xmlquery.Parse(body)
	if err != nil {
		return "", err
	}

	node, err := xmlquery.Query(root, "//NewExternalIPAddress/text()")
	if err != nil {
		return "", err
	}

	return node.Data, nil
}
