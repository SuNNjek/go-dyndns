package updater

import (
	"fmt"
	"github.com/google/wire"
	"github.com/kelseyhightower/envconfig"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var dynDnsSet = wire.NewSet(loadDynDnsConfig, newDynDnsUpdater)

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
	config *dynDnsUpdaterConfig
}

func newDynDnsUpdater(config *dynDnsUpdaterConfig) *dynDnsUpdater {
	return &dynDnsUpdater{config: config}
}

func (u *dynDnsUpdater) UpdateIP(addr net.IP) error {
	updateUrl, err := u.createUpdateUrl(addr)
	if err != nil {
		return err
	}

	resp, err := http.Get(updateUrl.String())
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	return nil
}

func (u *dynDnsUpdater) createUpdateUrl(addr net.IP) (*url.URL, error) {
	strUrl := fmt.Sprintf("https://%s/v3/update", u.config.Host)
	updateUrl, err := url.Parse(strUrl)
	if err != nil {
		return nil, err
	}

	password, err := u.getPassword()
	if err != nil {
		return nil, err
	}

	updateUrl.User = url.UserPassword(u.config.User, password)

	query := updateUrl.Query()
	query.Add("hostname", strings.Join(u.config.Domains, ","))
	query.Add("myip", fmt.Sprint(addr))
	updateUrl.RawQuery = query.Encode()

	return updateUrl, nil
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
