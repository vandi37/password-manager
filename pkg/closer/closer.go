package closer

import (
	"context"
	"sync"

	"github.com/vandi37/password-manager/pkg/logger"
	"github.com/vandi37/vanerrors"
)

const (
	ShutdownCancelled = "shutdown cancelled"
	GotSomeErrors     = "got some errors"
	ContextDone       = "context done"
)

type Fn func(ctx context.Context) error

func ToFn(f func() error) Fn {
	return func(ctx context.Context) error {
		err := make(chan (error))
		go func() {
			err <- f()
		}()
		select {
		case <-ctx.Done():
			return vanerrors.NewSimple(ContextDone)
		case res := <-err:
			return res
		}
	}
}

type Closer struct {
	logger *logger.Logger
	mu     sync.Mutex
	funcs  []Fn
}

func (c *Closer) Add(fn Fn) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs = append(c.funcs, fn)
}

func (c *Closer) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var (
		errs = make([]error, 0, len(c.funcs))
		wg   sync.WaitGroup
	)

	for _, f := range c.funcs {
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
		return vanerrors.NewWrap(ShutdownCancelled, ctx.Err(), vanerrors.EmptyHandler)
	}

	if len(errs) > 0 {
		for _, err := range errs {
			c.logger.Errorln(err)
		}
		return vanerrors.NewSimple(GotSomeErrors)
	}

	return nil
}

func New(logger *logger.Logger) *Closer {
	return &Closer{logger: logger}
}
