package addrproviders

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-dyndns/util"
	"io"
	"net/http"
	"testing"
)

func Test_webProvider_GetIP_Success(t *testing.T) {
	config := &WebProviderConfig{
		Url: "http://checkip.dyndns.com/",
	}

	httpClient := setupWebHttpClientMockWithIp("http://checkip.dyndns.com/", "127.0.0.1")

	provider := newWebProvider(config, httpClient)
	ip, err := provider.GetIP()

	assert.Nil(t, err)
	assert.Equal(t, "127.0.0.1", ip.String())

	httpClient.AssertExpectations(t)
}

func Test_webProvider_GetIP_MalformedIP(t *testing.T) {
	config := &WebProviderConfig{
		Url: "http://checkip.dyndns.com/",
	}

	httpClient := setupWebHttpClientMockWithIp("http://checkip.dyndns.com/", "asdf")

	provider := newWebProvider(config, httpClient)
	_, err := provider.GetIP()

	assert.ErrorIs(t, err, ParseIpError)
	httpClient.AssertExpectations(t)
}

func Test_webProvider_GetIP_InvalidResponseHtml(t *testing.T) {
	config := &WebProviderConfig{
		Url: "http://checkip.dyndns.com/",
	}

	httpClient := setupWebHttpClientMock("http://checkip.dyndns.com/", "asdf")

	provider := newWebProvider(config, httpClient)
	_, err := provider.GetIP()

	assert.ErrorIs(t, err, InvalidResponseError)
	httpClient.AssertExpectations(t)
}

func setupWebHttpClientMockWithIp(url, ip string) *util.MockHttpClient {
	return setupWebHttpClientMock(
		url,
		fmt.Sprintf("<html><head><title>Current IP Check</title></head><body>Current IP Address: %s</body></html>", ip),
	)
}

func setupWebHttpClientMock(url, responseText string) *util.MockHttpClient {
	httpClient := new(util.MockHttpClient)
	httpClient.On("Get", url).
		Return(&http.Response{
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
