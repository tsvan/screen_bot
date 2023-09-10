package internal

import "context"


type Script interface {
	Start(ctx context.Context) error
	Stop()
}
