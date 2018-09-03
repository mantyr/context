package context

import (
	"context"
)

// WGContext это интерфейс для WaitGroupContext
type WGContext interface {
	context.Context
	WContext

	Add(delta int) error
	Delete()
}

// WContext это интерфейс для контекстов с ожиданием
type WContext interface {
	Wait() <-chan struct{}
}
