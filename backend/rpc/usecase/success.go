package usecase

import (
	"context"
)

// Success :
type Success interface {
	GetSuccess(ctx context.Context) (result string, err error)
}
