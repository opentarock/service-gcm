package gcm

import (
	"time"

	"code.google.com/p/go.net/context"
	gcmlib "github.com/alexjlockwood/gcm"
	"github.com/cenkalti/backoff"

	"github.com/opentarock/service-api/go/util/contextutil"
)

const maxRetries = 10

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
	ticker, cancel := initBackoff()
	msg.DryRun = s.DryRun
	return contextutil.DoWithCancel(ctx, cancel, func() error {
		count := 0
		var err error
		for _ = range ticker.C {
			if count == maxRetries {
				return err
			}
			_, err = s.sender.SendNoRetry(msg)
			if err == nil {
				cancel()
				break
			}
			count++
		}
		return nil
	})
}

func initBackoff() (*backoff.Ticker, func()) {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 10 * time.Millisecond
	ticker := backoff.NewTicker(b)
	cancel := func() { ticker.Stop() }
	return ticker, cancel
}
