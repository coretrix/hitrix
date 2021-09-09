package fcm

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

type FCM interface {
	Send(ctx context.Context, message *messaging.Message) (string, error)
	SendDryRun(ctx context.Context, message *messaging.Message) (string, error)
	SendAll(ctx context.Context, messages []*messaging.Message) (*messaging.BatchResponse, error)
	SendAllDryRun(ctx context.Context, messages []*messaging.Message) (*messaging.BatchResponse, error)
	SendMulticast(ctx context.Context, message *messaging.MulticastMessage) (*messaging.BatchResponse, error)
	SendMulticastDryRun(ctx context.Context, message *messaging.MulticastMessage) (*messaging.BatchResponse, error)
	SubscribeToTopic(ctx context.Context, tokens []string, topic string) (*messaging.TopicManagementResponse, error)
	UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) (*messaging.TopicManagementResponse, error)
}

func NewFCM(ctx context.Context) (FCM, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, err
	}

	return app.Messaging(ctx)
}
