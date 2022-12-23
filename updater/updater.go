package updater

import (
	"context"
)

type Updater interface {
	UpdateIP(ctx context.Context, req *UpdateRequest) error
}
