package updater

import (
	"context"
	"net"
)

type Updater interface {
	UpdateIP(ctx context.Context, addr net.IP) error
}
