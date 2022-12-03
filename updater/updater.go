package updater

import (
	"net"
)

type Updater interface {
	UpdateIP(addr net.IP) error
}
