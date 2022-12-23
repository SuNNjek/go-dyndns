package updater

import (
	"bytes"
	"fmt"
	"go-dyndns/util"
	"net"
)

type UpdateRequest struct {
	IPv4       net.IP
	IPv6Prefix *util.IPv6Prefix
}

func (u *UpdateRequest) String() string {
	res := fmt.Sprintf("IPv4: %v", u.IPv4)

	if u.IPv6Prefix != nil {
		res += fmt.Sprintf("; IPv6 prefix: %v", u.IPv6Prefix)
	}

	return res
}

func (u *UpdateRequest) Equal(other *UpdateRequest) bool {
	if u == nil || other == nil {
		return false
	}

	// Either (both are nil) or (not nil and all the fields match)
	hasSamePrefix := (u.IPv6Prefix == nil && other.IPv6Prefix == nil) ||
		(u.IPv6Prefix != nil && other.IPv6Prefix != nil &&
			u.IPv6Prefix.IP.Equal(other.IPv6Prefix.IP) &&
			bytes.Equal(u.IPv6Prefix.Mask, other.IPv6Prefix.Mask))

	return u.IPv4.Equal(other.IPv4) && hasSamePrefix
}
