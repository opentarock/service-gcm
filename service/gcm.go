package service

import (
	"encoding/json"
	"log"
	"strings"

	gcmlib "github.com/alexjlockwood/gcm"

	"github.com/opentarock/service-api/go/proto"
	"github.com/opentarock/service-api/go/proto_errors"
	"github.com/opentarock/service-api/go/proto_gcm"
	"github.com/opentarock/service-api/go/service"
	"github.com/opentarock/service-gcm/gcm"
)

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
		var request proto_gcm.SendMessageRequest
		err := msg.Unmarshal(&request)
		if err != nil {
			log.Println(err)
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
		log.Printf("Sending message to %s", strings.Join(request.GetRegistrationIds(), ","))
		var data map[string]interface{}
		if request.GetData() != "" {
			err = json.Unmarshal([]byte(request.GetData()), &data)
			if err != nil {
				log.Println("Malformed JSON data: ", err)
				response.ErrorCode = proto_gcm.SendMessageResponse_MALFORMED_JSON.Enum()
				return proto.CompositeMessage{Message: &response}
			}
		}
		gcmMessage := &gcmlib.Message{
			RegistrationIDs: request.GetRegistrationIds(),
			Data:            data,
		}
		// Ignore the returned error because the only implementation always
		// returns nil.
		s.gcmSender.SendMessage(gcmMessage)

		return proto.CompositeMessage{Message: &response}
	})
}
