package context

import (
	"context"
	"os"
	"os/signal"
)

// waitSyscallContext это контекст который закрывается по сигналу
type waitSyscallContext struct {
	context.Context
	cancel context.CancelFunc
	sigs   chan os.Signal
}

// WaitSyscallContext возвращает контекст который закрывается по сигналу
//
// Внимание, перед использованием изучите работу signal.Notify
func WaitSyscallContext(
	ctx context.Context,
	signals ...os.Signal,
) (
	context.Context,
	context.CancelFunc,
) {
	c := &waitSyscallContext{}
	c.Context, c.cancel = context.WithCancel(ctx)
	c.sigs = make(chan os.Signal, 1)

	signal.Notify(
		c.sigs,
		signals...,
	)
	go c.run()
	return c, c.cancel
}

func (ctx *waitSyscallContext) run() {
	select {
	case <-ctx.sigs:
		ctx.cancel()
	case <-ctx.Done():
		// exit
	}
}
