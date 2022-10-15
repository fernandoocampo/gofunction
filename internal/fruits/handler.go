package fruits

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func MakeHandler(service *Service) func(context.Context, events.SQSEvent) error {
	return func(ctx context.Context, sqsEvent events.SQSEvent) error {
		for _, message := range sqsEvent.Records {
			log.Println("new event received", "id", message.MessageId, "source", message.EventSource, "content", message.Body)

			service.Audit(ctx, message.Body)
		}

		return nil
	}
}
