package gcm

import (
	"code.google.com/p/go.net/context"
	gcmlib "github.com/alexjlockwood/gcm"
)

type Sender interface {
	SendMessage(ctx context.Context, msg *gcmlib.Message) error
}
