package context

import (
	"context"
	"syscall"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWaitSyscallContext(t *testing.T) {
	Convey("Проверяем WaitSyscallContext", t, func() {
		ctx, cancel := WaitSyscallContext(
			context.Background(),
			syscall.SIGINT,
			syscall.SIGTERM,
		)
		So(ctx, ShouldNotBeNil)
		So(cancel, ShouldNotBeNil)
		So(ctx, ShouldNotBeDone)
		Convey("Проверяем закрытие контекста", func() {
			Convey("CancelFunc", func() {
				cancel()
				So(ctx, ShouldBeDone)
			})
			Convey("syscall.SIGINT", func() {
				err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
				So(err, ShouldBeNil)
				time.Sleep(1 * time.Second)
				So(ctx, ShouldBeDone)
			})
			Convey("syscall.SIGTERM", func() {
				err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
				So(err, ShouldBeNil)
				time.Sleep(1 * time.Second)
				So(ctx, ShouldBeDone)
			})
		})
		Convey("Проверяем что контекст нельзя закрыть", func() {
			Convey("syscall.SIGUSR1", func() {
				err := syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
				So(err, ShouldBeNil)
				time.Sleep(1 * time.Second)
				So(ctx, ShouldNotBeDone)
			})
		})
	})
}
