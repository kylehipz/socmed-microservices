package events

import "context"

type Consumer interface {
	Start(context.Context)
	Stop(context.Context)
}
