package context

import (
	"context"
	"errors"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func ShouldBeDone(ctx context.Context) {
	var err error
	select {
	case <-ctx.Done():
	default:
		err = errors.New("context not done")
	}
	So(err, ShouldBeNil)
}

func ShouldNotBeDone(ctx context.Context) {
	var err error
	select {
	case <-ctx.Done():
		err = errors.New("context done")
	default:
	}
	So(err, ShouldBeNil)
}

func ShouldBeWait(ctx WContext) {
	var err error
	select {
	case <-ctx.Wait():
	default:
		err = errors.New("context not wait")
	}
	So(err, ShouldBeNil)
}

func ShouldNotBeWait(ctx WContext) {
	var err error
	select {
	case <-ctx.Wait():
		err = errors.New("context wait")
	default:
	}
	So(err, ShouldBeNil)
}

func TestWaitGroupContext(t *testing.T) {
	Convey("Проверяем WaitGroupContext", t, func() {
		Convey("Проверяем закрытие контекста", func() {
			ctx, cancel := WaitGroupContext(context.Background())
			Convey("Добавили элементы", func() {
				err := ctx.Add(2)
				So(err, ShouldBeNil)
				ShouldNotBeDone(ctx)
				ShouldNotBeWait(ctx)
				Convey("Закрыли контекст", func() {
					cancel()
					time.Sleep(10 * time.Millisecond)
					ShouldBeDone(ctx)
					ShouldNotBeWait(ctx)
					Convey("Не можем добавить элементы так как контекст закрыт", func() {
						err := ctx.Add(1)
						So(err, ShouldNotBeNil)
					})
					Convey("Убрали элемент", func() {
						ctx.Delete()
						ShouldBeDone(ctx)
						ShouldNotBeWait(ctx)
						Convey("Убираем все элементы - wait закрылся", func() {
							ctx.Delete()
							ShouldBeDone(ctx)
							ShouldBeWait(ctx)
						})
						Convey("Убрали на несколько элементов больше - wait закрылся", func() {
							ctx.Delete()
							ctx.Delete()
							ctx.Delete()
							ShouldBeDone(ctx)
							ShouldBeWait(ctx)
						})
					})
				})
			})
			Convey("Закрытие от родительского контекста", func() {
				parentContext, parentCancel := context.WithCancel(context.Background())
				ctx, cancel := WaitGroupContext(parentContext)

				Convey("Родительский контекст не закрыт", func() {
					Convey("Done", func() {
						ShouldNotBeDone(ctx)
					})
					Convey("Wait", func() {
						ShouldNotBeWait(ctx)
					})
				})
				Convey("Родительский контекст закрыт", func() {
					parentCancel()
					time.Sleep(10 * time.Millisecond)
					Convey("Done", func() {
						ShouldBeDone(ctx)
					})
					Convey("Wait", func() {
						ShouldBeWait(ctx)
					})
				})
				Convey("Принудительно закрыли самостоятельно", func() {
					cancel()
					time.Sleep(10 * time.Millisecond)
					Convey("Done", func() {
						ShouldBeDone(ctx)
					})
					Convey("Wait", func() {
						ShouldBeWait(ctx)
					})
				})
			})
		})
	})
}
