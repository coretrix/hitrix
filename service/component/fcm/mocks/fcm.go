package mocks

import (
	"context"

	"firebase.google.com/go/messaging"
	"github.com/stretchr/testify/mock"
)

type FakeFCM struct {
	mock.Mock
}

func (f *FakeFCM) Send(_ context.Context, message *messaging.Message) (string, error) {
	args := f.Called(message)

	return args.String(0), args.Error(1)
}

func (f *FakeFCM) SendDryRun(_ context.Context, message *messaging.Message) (string, error) {
	args := f.Called(message)

	return args.String(0), args.Error(1)
}

func (f *FakeFCM) SendAll(_ context.Context, messages []*messaging.Message) (*messaging.BatchResponse, error) {
	args := f.Called(messages)

	return args.Get(0).(*messaging.BatchResponse), args.Error(1)
}

func (f *FakeFCM) SendAllDryRun(_ context.Context, messages []*messaging.Message) (*messaging.BatchResponse, error) {
	args := f.Called(messages)

	return args.Get(0).(*messaging.BatchResponse), args.Error(1)
}

func (f *FakeFCM) SendMulticast(_ context.Context, message *messaging.MulticastMessage) (*messaging.BatchResponse, error) {
	args := f.Called(message)

	return args.Get(0).(*messaging.BatchResponse), args.Error(1)
}

func (f *FakeFCM) SendMulticastDryRun(_ context.Context, message *messaging.MulticastMessage) (*messaging.BatchResponse, error) {
	args := f.Called(message)

	return args.Get(0).(*messaging.BatchResponse), args.Error(1)
}

func (f *FakeFCM) SubscribeToTopic(_ context.Context, tokens []string, topic string) (*messaging.TopicManagementResponse, error) {
	args := f.Called(tokens, topic)

	return args.Get(0).(*messaging.TopicManagementResponse), args.Error(1)
}

func (f *FakeFCM) UnsubscribeFromTopic(_ context.Context, tokens []string, topic string) (*messaging.TopicManagementResponse, error) {
	args := f.Called(tokens, topic)

	return args.Get(0).(*messaging.TopicManagementResponse), args.Error(1)
}
