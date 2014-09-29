package gcm

import gcmlib "github.com/alexjlockwood/gcm"

type Sender interface {
	SendMessage(msg *gcmlib.Message) error
}
