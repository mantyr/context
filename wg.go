package context

import (
	"context"
	"errors"
	"sync"
)

var (
	// closedchan это многоразовый закрытый канал
	closedchan = make(chan struct{})

	// canceled это ошибка свидетельствует о том что контекст закрыт
	canceled = CanceledError("context canceled")
)

// CanceledError эта ошибка говорит о том что контекст закрыт
type CanceledError string

// Error реализует интерфейс error
func (c CanceledError) Error() string {
	return string(c)
}

type waitGroupContext struct {
	context.Context

	mu         sync.Mutex
	group      int
	done       chan struct{}
	wait       chan struct{}
	err        error
	waitClosed bool
}

// WaitGroupContext возвращает контекст с функциями WaitGroup
// 1. Можно добавить элемент только если контекст не закрыт
// 2. Можно убрать элемент в любой момент
// 3. Можно закрыть контекст в любой момент
// 4. Можно дождаться закрытия контекста
// 5. Можно дождаться завершения операций связанных с контекстом
func WaitGroupContext(
	parent context.Context,
) (
	context.Context,
	context.CancelFunc,
) {
	ctx := &waitGroupContext{
		Context: parent,
	}
	if parent != nil {
		go func() {
			select {
			case <-parent.Done():
				ctx.cancel(parent.Err())
			case <-ctx.Done():
			}
		}()
	}
	return ctx, func() { ctx.cancel(canceled) }
}

func (c *waitGroupContext) cancel(err error) {
	if err == nil {
		err = canceled
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.err != nil {
		return
	}

	c.err = err
	if c.done == nil {
		c.done = closedchan
	} else {
		close(c.done)
	}
	if c.group > 0 {
		return
	}
	c.waitClosed = true
	if c.wait == nil {
		c.wait = closedchan
	} else {
		close(c.wait)
	}
}

// Err возвращает nil если контекст не закрыт и ошибку если закрыт
func (c *waitGroupContext) Err() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	err := c.err
	return err
}

// Add добавляет элемент в контекст только если контекст не закрыт
func (c *waitGroupContext) Add(delta int) error {
	if delta < 1 {
		return errors.New("expected delta > 0")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err != nil {
		return c.err
	}
	c.group += delta
	return nil
}

// Delete убирает элемент из контекста
func (c *waitGroupContext) Delete() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.group == 0 {
		return
	}
	c.group--
	switch {
	case c.group > 0:
		return
	case c.err == nil:
		return
	case c.waitClosed:
		return
	}
	if c.wait == nil {
		c.wait = closedchan
	} else {
		close(c.wait)
	}
}

// Done возвращает канал для ожидания закрытия контекста
// nolint: dupl
func (c *waitGroupContext) Done() <-chan struct{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.done == nil {
		c.done = make(chan struct{})
	}
	d := c.done
	return d
}

// Wait возвращает канал для ожидания закрытия всех дочерних задач
// nolint: dupl
func (c *waitGroupContext) Wait() <-chan struct{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.wait == nil {
		c.wait = make(chan struct{})
	}
	w := c.wait
	return w
}
