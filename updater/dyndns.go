package updater

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/wire"
	"github.com/kelseyhightower/envconfig"
	"go-dyndns/log"
	"go-dyndns/util"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var dynDnsSet = wire.NewSet(loadDynDnsConfig, newDynDnsUpdater, wire.FieldsOf(new(*dynDnsUpdaterConfig), "PasswordFile"))

var (
	AuthenticationError = errors.New("the given username or password are incorrect")
	HostNotFQDN         = errors.New("the given host is not a fully qualified domain name")
	HostNotFound        = errors.New("the given host is not associated with the given user")
	TooManyHosts        = errors.New("too many hosts given")
	AbuseError          = errors.New("IP was updated too often and has been blocked for update abuse")
	DnsError            = errors.New("DNS error")
	ServerError         = errors.New("server error")
)

type InvalidResponse struct {
	response string
}

func (i *InvalidResponse) Error() string {
	return fmt.Sprintf("invalid response: %s", i.response)
}

type dynDnsUpdaterConfig struct {
	Host         string                `required:"true"`
	User         string                `required:"true"`
	PasswordFile util.PasswordFilePath `required:"true"`
	Domains      []string              `required:"true"`
}

func loadDynDnsConfig() (*dynDnsUpdaterConfig, error) {
	var config dynDnsUpdaterConfig
	if err := envconfig.Process("dyndns", &config); err != nil {
		return nil, err
	}

	return &config, nil
}

type dynDnsUpdater struct {
	config           *dynDnsUpdaterConfig
	passwordProvider util.PasswordProvider
	logger           log.Logger
	httpClient       util.HttpClient
}

func newDynDnsUpdater(
	config *dynDnsUpdaterConfig,
	passwordProvider util.PasswordProvider,
	logger log.Logger,
	httpClient util.HttpClient,
) *dynDnsUpdater {
	return &dynDnsUpdater{
		config:           config,
		passwordProvider: passwordProvider,
		logger:           logger,
		httpClient:       httpClient,
	}
}

func (u *dynDnsUpdater) UpdateIP(ctx context.Context, req *UpdateRequest) error {
	u.logger.Info("Updating IP for domains %v to %v", strings.Join(u.config.Domains, ", "), req)

	webReq, err := u.createUpdateWebRequest(ctx, req)
	if err != nil {
		return err
	}

	resp, err := u.httpClient.Do(webReq)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return u.handleResponse(strings.TrimSpace(string(body)))
}

func (u *dynDnsUpdater) createUpdateWebRequest(ctx context.Context, req *UpdateRequest) (*http.Request, error) {
	strUrl := fmt.Sprintf("https://%s/v3/update", u.config.Host)
	updateUrl, err := url.Parse(strUrl)
	if err != nil {
		return nil, err
	}

	query := make(url.Values)
	query.Add("hostname", strings.Join(u.config.Domains, ","))
	query.Add("myip", req.IPv4.String())

	if req.IPv6Prefix != nil {
		query.Add("ip6lanprefix", req.IPv6Prefix.String())
	}

	updateUrl.RawQuery = query.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, updateUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	password, err := u.passwordProvider.GetPassword()
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(u.config.User, password)

	u.logger.Trace(
		"Created update request for URL %v with user %s and password %s",
		updateUrl.String(), u.config.User, password,
	)

	return request, nil
}

func (u *dynDnsUpdater) handleResponse(response string) error {
	u.logger.Trace("Received response: %s", response)

	if strings.HasPrefix(response, "good") {
		return nil
	}

	if strings.HasPrefix(response, "nochg") {
		return nil
	}

	if strings.HasPrefix(response, "badauth") {
		return AuthenticationError
	}

	if strings.HasPrefix(response, "notfqdn") {
		return HostNotFQDN
	}

	if strings.HasPrefix(response, "nohost") {
		return HostNotFound
	}

	if strings.HasPrefix(response, "numhost") {
		return TooManyHosts
	}

	if strings.HasPrefix(response, "abuse") || strings.HasPrefix(response, "toomanyrequests") {
		return AbuseError
	}

	if strings.HasPrefix(response, "dnserr") {
		return DnsError
	}

	if strings.HasPrefix(response, "911") {
		return ServerError
	}

	return &InvalidResponse{response: response}
}
