package closer

import (
	"context"
	"go.uber.org/zap"
	"sync"

	"github.com/vandi37/password-manager/pkg/logger"
	"github.com/vandi37/vanerrors"
)

const (
	ShutdownCancelled = "shutdown cancelled"
	GotSomeErrors     = "got some errors"
)

type Fn func(ctx context.Context) error

type Closer struct {
	mu  sync.Mutex
	fns []Fn
}

func (c *Closer) Add(fn Fn) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.fns = append(c.fns, fn)
}

func (c *Closer) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var (
		errs = make([]error, 0, len(c.fns))
		wg   sync.WaitGroup
	)

	for _, f := range c.fns {
		wg.Add(1)
		go func(f Fn) {
			defer wg.Done()

			if err := f(ctx); err != nil {
				errs = append(errs, err)
			}

		}(f)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		break
	case <-ctx.Done():
		return vanerrors.Simple(ShutdownCancelled)
	}

	if len(errs) > 0 {
		for _, err := range errs {
			logger.Error(ctx, "Close error", zap.Error(err))
		}
		return vanerrors.Simple(GotSomeErrors)
	}

	return nil
}

func New() *Closer {
	return &Closer{}
}
