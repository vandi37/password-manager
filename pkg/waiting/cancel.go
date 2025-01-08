package waiting

type Cancel struct {
	cancel chan struct{}
}

func NewCancel() Cancel {
	return Cancel{cancel: make(chan struct{})}
}

func (c Cancel) Cancel() {
	c.cancel <- struct{}{}
}

func (c Cancel) Canceled() chan struct{} {
	return c.cancel
}
