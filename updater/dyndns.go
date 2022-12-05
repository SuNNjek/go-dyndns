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
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var dynDnsSet = wire.NewSet(loadDynDnsConfig, newDynDnsUpdater)

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
	Host         string   `required:"true"`
	User         string   `required:"true"`
	PasswordFile string   `required:"true"`
	Domains      []string `required:"true"`
}

func loadDynDnsConfig() (*dynDnsUpdaterConfig, error) {
	var config dynDnsUpdaterConfig
	if err := envconfig.Process("dyndns", &config); err != nil {
		return nil, err
	}

	return &config, nil
}

type dynDnsUpdater struct {
	config     *dynDnsUpdaterConfig
	logger     log.Logger
	httpClient util.HttpClient
}

func newDynDnsUpdater(config *dynDnsUpdaterConfig, logger log.Logger, httpClient util.HttpClient) *dynDnsUpdater {
	return &dynDnsUpdater{config: config, logger: logger, httpClient: httpClient}
}

func (u *dynDnsUpdater) UpdateIP(ctx context.Context, addr net.IP) error {
	u.logger.Info("Updating IP for domains %v to %v", strings.Join(u.config.Domains, ", "), addr)

	req, err := u.createUpdateRequest(ctx, addr)
	if err != nil {
		return err
	}

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return handleResponse(strings.TrimSpace(string(body)))
}

func (u *dynDnsUpdater) createUpdateRequest(ctx context.Context, addr net.IP) (*http.Request, error) {
	strUrl := fmt.Sprintf("https://%s/v3/update", u.config.Host)
	updateUrl, err := url.Parse(strUrl)
	if err != nil {
		return nil, err
	}

	query := make(url.Values)
	query.Add("hostname", strings.Join(u.config.Domains, ","))
	query.Add("myip", fmt.Sprint(addr))
	updateUrl.RawQuery = query.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, updateUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	password, err := u.getPassword()
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

func handleResponse(response string) error {
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

func (u *dynDnsUpdater) getPassword() (string, error) {
	file, err := os.Open(u.config.PasswordFile)
	if err != nil {
		return "", err
	}

	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
