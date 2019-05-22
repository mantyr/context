package context

import (
	"context"
	"errors"

	"github.com/smartystreets/goconvey/convey"
)

// ShouldBeDone проверяет что контекст закрыт
func ShouldBeDone(actual interface{}, expected ...interface{}) string {
	var err error
	ctx, ok := actual.(context.Context)
	if !ok {
		err = errors.New("actual is not context")
	}
	select {
	case <-ctx.Done():
	default:
		err = errors.New("context not done")
	}
	return convey.ShouldBeNil(err)
}

// ShouldNotBeDone проверяет что контекст не закрыт
func ShouldNotBeDone(actual interface{}, expected ...interface{}) string {
	var err error
	ctx, ok := actual.(context.Context)
	if !ok {
		err = errors.New("actual is not context")
	}
	select {
	case <-ctx.Done():
		err = errors.New("context done")
	default:
	}
	return convey.ShouldBeNil(err)
}

// ShouldBeWait проверяет что все задачи в контексте были завершены
func ShouldBeWait(actual interface{}, expected ...interface{}) string {
	var err error
	ctx, ok := actual.(WContext)
	if !ok {
		err = errors.New("actual is not WContext")
	}
	select {
	case <-ctx.Wait():
	default:
		err = errors.New("context not wait")
	}
	return convey.ShouldBeNil(err)
}

// ShouldNotBeWait проверяет что в контексте есть не завершённые задачи
func ShouldNotBeWait(actual interface{}, expected ...interface{}) string {
	var err error
	ctx, ok := actual.(WContext)
	if !ok {
		err = errors.New("actual is not WContext")
	}
	select {
	case <-ctx.Wait():
		err = errors.New("context wait")
	default:
	}
	return convey.ShouldBeNil(err)
}
