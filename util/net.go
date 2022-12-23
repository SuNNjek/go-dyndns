package util

import (
	"net"
)

// IPv6Prefix is a wrapper around net.IPNet with better text (un)marshalling
type IPv6Prefix struct {
	net.IPNet
}

func (p *IPv6Prefix) UnmarshalText(text []byte) error {
	_, ipNet, err := net.ParseCIDR(string(text))
	if err != nil {
		return err
	}

	p.IPNet = *ipNet
	return nil
}

func (p *IPv6Prefix) MarshalText() (text []byte, err error) {
	return []byte(p.String()), nil
}

func ParseIPv6Prefix(str string) (*IPv6Prefix, error) {
	_, ipNet, err := net.ParseCIDR(str)
	if err != nil {
		return nil, err
	}

	return &IPv6Prefix{*ipNet}, nil
}
