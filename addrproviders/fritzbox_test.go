package addrproviders

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-dyndns/util"
	"io"
	"net/http"
	"testing"
)

func Test_fritzBoxProvider_GetIP_Success(t *testing.T) {
	config := &fritzBoxConfig{
		Host: "fritz.box",
	}

	httpClient := setupFritzboxHttpClientMockWithIP("127.0.0.1")

	provider := newFritzBoxProvider(config, httpClient)
	ip, err := provider.GetIP()

	assert.Nil(t, err)
	assert.Equal(t, "127.0.0.1", ip.String())

	httpClient.AssertExpectations(t)
}

func Test_fritzBoxProvider_GetIP_MalformedIP(t *testing.T) {
	config := &fritzBoxConfig{
		Host: "fritz.box",
	}

	httpClient := setupFritzboxHttpClientMockWithIP("asdf")

	provider := newFritzBoxProvider(config, httpClient)
	_, err := provider.GetIP()

	assert.ErrorIs(t, err, ParseIpError)
	httpClient.AssertExpectations(t)
}

func Test_fritzBoxProvider_GetIP_InvalidResponseXml(t *testing.T) {
	config := &fritzBoxConfig{
		Host: "fritz.box",
	}

	httpClient := setupFritzboxHttpClientMock("asdf")

	provider := newFritzBoxProvider(config, httpClient)
	_, err := provider.GetIP()

	assert.ErrorIs(t, err, InvalidResponseError)
	httpClient.AssertExpectations(t)
}

func setupFritzboxHttpClientMockWithIP(ip string) *util.MockHttpClient {
	responseText := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
<s:Body>
<u:GetExternalIPAddressResponse xmlns:u="urn:schemas-upnp-org:service:WANIPConnection:1">
<NewExternalIPAddress>%s</NewExternalIPAddress>
</u:GetExternalIPAddressResponse>
</s:Body>
</s:Envelope>`, ip)

	return setupFritzboxHttpClientMock(responseText)
}

func setupFritzboxHttpClientMock(responseText string) *util.MockHttpClient {
	httpClient := new(util.MockHttpClient)
	httpClient.On(
		"Do",
		mock.MatchedBy(func(req *http.Request) bool {
			body, err := req.GetBody()
			if err != nil {
				return false
			}

			defer body.Close()

			text, err := io.ReadAll(body)
			return string(text) == `<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
	<s:Body>
		<u:GetExternalIPAddress xmlns:u="urn:schemas-upnp-org:service:WANIPConnection:1" />
	</s:Body>
</s:Envelope>`
		}),
	).Return(&http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          io.NopCloser(bytes.NewBufferString(responseText)),
		ContentLength: int64(len(responseText)),
		Header:        make(http.Header, 0),
	}, nil)

	return httpClient
}
