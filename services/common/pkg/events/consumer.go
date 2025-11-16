package events

import "context"

type Consumer interface {
	Start(context.Context) error
	Wait(context.Context)
}
