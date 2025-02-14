package server

import (
	"context"
)

type Server interface {
	ReadLoop()
	Shutdown(ctx context.Context)
}
