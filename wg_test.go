package context

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWaitGroupContext(t *testing.T) {
	Convey("Проверяем WaitGroupContext", t, func() {
		Convey("Проверяем закрытие контекста", func() {
			ctx, cancel := WaitGroupContext(context.Background())
			Convey("Добавили элементы", func() {
				err := ctx.Add(2)
				So(err, ShouldBeNil)
				So(ctx, ShouldNotBeDone)
				So(ctx, ShouldNotBeWait)
				Convey("Закрыли контекст", func() {
					cancel()
					time.Sleep(10 * time.Millisecond)
					So(ctx, ShouldBeDone)
					So(ctx, ShouldNotBeWait)
					Convey("Не можем добавить элементы так как контекст закрыт", func() {
						err := ctx.Add(1)
						So(err, ShouldNotBeNil)
					})
					Convey("Убрали элемент", func() {
						ctx.Delete()
						So(ctx, ShouldBeDone)
						So(ctx, ShouldNotBeWait)
						Convey("Убираем все элементы - wait закрылся", func() {
							ctx.Delete()
							So(ctx, ShouldBeDone)
							So(ctx, ShouldBeWait)
						})
						Convey("Убрали на несколько элементов больше - wait закрылся", func() {
							ctx.Delete()
							ctx.Delete()
							ctx.Delete()
							So(ctx, ShouldBeDone)
							So(ctx, ShouldBeWait)
						})
					})
				})
			})
			Convey("Закрытие от родительского контекста", func() {
				parentContext, parentCancel := context.WithCancel(context.Background())
				ctx, cancel := WaitGroupContext(parentContext)

				Convey("Родительский контекст не закрыт", func() {
					Convey("Done", func() {
						So(ctx, ShouldNotBeDone)
					})
					Convey("Wait", func() {
						So(ctx, ShouldNotBeWait)
					})
				})
				Convey("Родительский контекст закрыт", func() {
					parentCancel()
					time.Sleep(10 * time.Millisecond)
					Convey("Done", func() {
						So(ctx, ShouldBeDone)
					})
					Convey("Wait", func() {
						So(ctx, ShouldBeWait)
					})
				})
				Convey("Принудительно закрыли самостоятельно", func() {
					cancel()
					time.Sleep(10 * time.Millisecond)
					Convey("Done", func() {
						So(ctx, ShouldBeDone)
					})
					Convey("Wait", func() {
						So(ctx, ShouldBeWait)
					})
				})
			})
		})
	})
}
