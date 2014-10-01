package gcm

import (
	"code.google.com/p/go.net/context"

	gcmlib "github.com/alexjlockwood/gcm"
	"github.com/opentarock/service-api/go/util/contextutil"
)

type RetrySender struct {
	sender *gcmlib.Sender
	DryRun bool
}

func NewRetrySender(apiKey string) *RetrySender {
	return &RetrySender{
		sender: &gcmlib.Sender{ApiKey: apiKey},
		DryRun: false,
	}
}

func (s *RetrySender) SendMessage(ctx context.Context, msg *gcmlib.Message) error {
	return contextutil.Do(ctx, func() error {
		_, err := s.sender.SendNoRetry(msg)
		return err
	})
}
