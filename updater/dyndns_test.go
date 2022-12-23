package updater

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-dyndns/log"
	"go-dyndns/util"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func Test_dynDnsUpdater_UpdateIP_shouldCallApiCorrectly(t *testing.T) {
	responseText := "good 127.0.0.1"

	httpClient := new(util.MockHttpClient)
	httpClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		correctUrl, _ := url.Parse("https://ddns.example.com/v3/update")
		query := make(url.Values)
		query.Set("hostname", strings.Join([]string{"example1.com", "example2.com"}, ","))
		query.Set("myip", "127.0.0.1")
		correctUrl.RawQuery = query.Encode()

		isUrlCorrect := req.URL.String() == correctUrl.String()

		user, password, hasAuth := req.BasicAuth()
		authCorrect := hasAuth && user == "user" && password == "password"

		return isUrlCorrect && authCorrect
	})).Return(&http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          io.NopCloser(bytes.NewBufferString(responseText)),
		ContentLength: int64(len(responseText)),
		Header:        make(http.Header, 0),
	}, nil)

	logger := log.CreateTestLogger()
	passwordProvider := &util.MockPasswordProvider{}
	passwordProvider.On("GetPassword").Return("password", nil)

	updater := newDynDnsUpdater(&dynDnsUpdaterConfig{
		Host:    "ddns.example.com",
		User:    "user",
		Domains: []string{"example1.com", "example2.com"},
	}, passwordProvider, logger, httpClient)

	err := updater.UpdateIP(context.Background(), &UpdateRequest{IPv4: net.ParseIP("127.0.0.1")})
	assert.Nil(t, err)

	httpClient.AssertExpectations(t)
	passwordProvider.AssertExpectations(t)
}

func isError(target error) assert.ErrorAssertionFunc {
	return func(t assert.TestingT, err error, i ...interface{}) bool {
		return assert.ErrorIs(t, err, target, i)
	}
}

func Test_dynDnsUpdater_handleResponse(t *testing.T) {
	type args struct {
		response string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Handle good",
			args:    args{response: "good 127.0.0.1"},
			wantErr: assert.NoError,
		},
		{
			name:    "Handle nochg",
			args:    args{response: "nochg 127.0.0.1"},
			wantErr: assert.NoError,
		},
		{
			name:    "Handle badauth",
			args:    args{response: "badauth"},
			wantErr: isError(AuthenticationError),
		},
		{
			name:    "Handle notfqdn",
			args:    args{response: "notfqdn"},
			wantErr: isError(HostNotFQDN),
		},
		{
			name:    "Handle nohost",
			args:    args{response: "nohost"},
			wantErr: isError(HostNotFound),
		},
		{
			name:    "Handle numhost",
			args:    args{response: "numhost"},
			wantErr: isError(TooManyHosts),
		},
		{
			name:    "Handle abuse",
			args:    args{response: "abuse"},
			wantErr: isError(AbuseError),
		},
		{
			name:    "Handle toomanyrequests",
			args:    args{response: "toomanyrequests"},
			wantErr: isError(AbuseError),
		},
		{
			name:    "Handle dnserr",
			args:    args{response: "dnserr"},
			wantErr: isError(DnsError),
		},
		{
			name:    "Handle 911",
			args:    args{response: "911"},
			wantErr: isError(ServerError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &dynDnsUpdater{
				config:           &dynDnsUpdaterConfig{},
				passwordProvider: &util.MockPasswordProvider{},
				logger:           log.CreateTestLogger(),
				httpClient:       &util.MockHttpClient{},
			}
			tt.wantErr(t, u.handleResponse(tt.args.response), fmt.Sprintf("handleResponse(%v)", tt.args.response))
		})
	}
}
