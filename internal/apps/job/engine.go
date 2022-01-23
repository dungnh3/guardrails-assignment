package job

import "context"

type Engine interface {
	Run(ctx context.Context) error
	Close(ctx context.Context) error
}
