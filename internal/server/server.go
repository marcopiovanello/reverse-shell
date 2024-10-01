package server

import (
	"context"
)

type Server interface {
	ReadLoop(ctx context.Context)
}
