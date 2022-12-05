# go-dyndns

[![codecov](https://codecov.io/gh/SuNNjek/go-dyndns/branch/main/graph/badge.svg?token=DWJII5ZVJ5)](https://codecov.io/gh/SuNNjek/go-dyndns)

go-dyndns is a DynDns client written in Golang (yes, very creative name, I know :smile: )

## Installation/Configuration

go-dyndns is designed to be run inside a Docker container. Therefore, configuration is done
via environment variables. The following is an example Docker compose file:

```yaml
version: '3.6'

secrets:
  password:
    file: password.txt

services:
  go-dyndns:
    image: ghcr.io/sunnjek/go-dyndns
    secrets:
      - password
    tmpfs: /tmp
    environment:
      # To instead retrieve the public IP from your FritzBox, use this instead of IPCHECK_URL:
      # CLIENT_IPPROVIDER: fritzbox
      IPCHECK_URL: 'http://checkip.dyndns.com'
      CACHE_FILE: /tmp/go-dyndns.cache
      DYNDNS_DOMAINS: example1.com,example2.com
      DYNDNS_HOST: members.dyndns.org
      DYNDNS_PASSWORDFILE: /run/secrets/password
      DYNDNS_USER: user
```

The following environment variables can be used:

* `CLIENT_IPPROVIDER`: Can be either `web` or `fritzbox`, default is `web`.
  Controls how the public IP is determined.
* `CLIENT_DELAY`: Delay between updates. Defaults to 10 minutes (`10m`)
* `FRITZBOX_HOST`: The hostname of the FritzBox. Defaults to `fritz.box`. Only used with `fritzbox` provider
* `IPCHECK_URL`: The URL to use to determine the public IP address. Only used with `web` provider
* `CACHE_FILE`: Specifies the location for the cache file. If not specified, the cache
  will be stored in memory
* `DYNDNS_DOMAINS`: Comma-separated list of domains to update
* `DYNDNS_HOST`: Host to which to send the update request.
* `DYNDNS_USER`: The username with which to authenticate
* `DYNDNS_PASSWORDFILE`: Path to a file containing the password with which to authenticate.
  Use a docker secret for this one (see example above)
* `LOG_LEVEL`: The level used for logging. Potentially useful for diagnosing problems
  (if I can be bothered to add more logs, that is :smile: )