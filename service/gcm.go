package service

import (
	"encoding/json"
	"strings"
	"time"

	"code.google.com/p/go.net/context"

	gcmlib "github.com/alexjlockwood/gcm"

	"github.com/opentarock/service-api/go/proto"
	"github.com/opentarock/service-api/go/proto_errors"
	"github.com/opentarock/service-api/go/proto_gcm"
	"github.com/opentarock/service-api/go/reqcontext"
	"github.com/opentarock/service-api/go/service"
	"github.com/opentarock/service-gcm/gcm"
)

const defaultTimeout = 1 * time.Second

type gcmServiceHandlers struct {
	gcmSender gcm.Sender
}

func NewGcmServiceHandlers(gcmSender gcm.Sender) *gcmServiceHandlers {
	return &gcmServiceHandlers{
		gcmSender: gcmSender,
	}
}

func (s *gcmServiceHandlers) SendMessageHandler() service.MessageHandler {
	return service.MessageHandlerFunc(func(msg *proto.Message) proto.CompositeMessage {
		ctx, cancel := context.WithTimeout(reqcontext.NewContext(context.Background(), msg), defaultTimeout)
		defer cancel()

		logger := reqcontext.ContextLogger(ctx)

		var request proto_gcm.SendMessageRequest
		err := msg.Unmarshal(&request)
		if err != nil {
			logger.Println(err)
			return proto.CompositeMessage{
				Message: proto_errors.NewMalformedMessage(request.GetMessageType()),
			}
		}

		if len(request.GetRegistrationIds()) < 1 {
			return proto.CompositeMessage{
				Message: proto_errors.NewMissingFieldError("registration_ids"),
			}
		}

		response := proto_gcm.SendMessageResponse{}
		logger.Printf("[INFO] Sending message to %s", strings.Join(request.GetRegistrationIds(), ","))
		var data map[string]interface{}
		if request.GetData() != "" {
			err = json.Unmarshal([]byte(request.GetData()), &data)
			if err != nil {
				logger.Println("[ERROR] Malformed JSON data: ", err)
				response.ErrorCode = proto_gcm.SendMessageResponse_MALFORMED_JSON.Enum()
				return proto.CompositeMessage{Message: &response}
			}
		}
		gcmMessage := &gcmlib.Message{
			RegistrationIDs: request.GetRegistrationIds(),
			Data:            data,
		}
		err = s.gcmSender.SendMessage(ctx, gcmMessage)
		if err != nil {
			logger.Println(err)
			return proto.CompositeMessage{
				Message: proto_errors.NewInternalError("Failed to send message."),
			}
		}

		return proto.CompositeMessage{Message: &response}
	})
}
