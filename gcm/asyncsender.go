package gcm

import (
	"log"

	gcmlib "github.com/alexjlockwood/gcm"
)

type AsyncSender struct {
	sender     *gcmlib.Sender
	numRetries uint
	DryRun     bool
}

func NewAsyncSender(apiKey string, numRetries uint) *AsyncSender {
	return &AsyncSender{
		sender:     &gcmlib.Sender{ApiKey: apiKey},
		numRetries: numRetries,
		DryRun:     false,
	}
}

func (s *AsyncSender) SendMessage(msg *gcmlib.Message) error {
	go func() {
		_, err := s.sender.Send(msg, int(s.numRetries))
		if err != nil {
			log.Printf("Failed to send message to '%s' after %d retries: %s",
				msg.RegistrationIDs, s.numRetries, err)
		}
	}()
	return nil
}
