package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type FakeKubernetes struct {
	mock.Mock
}

func (k *FakeKubernetes) GetIngressDomains(ctx context.Context) ([]string, error) {
	args := k.Called(ctx)

	return args.Get(0).([]string), args.Error(1)
}

func (k *FakeKubernetes) AddIngress(ctx context.Context, domain, secretName, serviceName, servicePortName string, annotations map[string]string) error {
	args := k.Called(ctx, domain, secretName, serviceName, servicePortName)

	return args.Error(0)
}

func (k *FakeKubernetes) RemoveIngress(ctx context.Context, domain string) error {
	args := k.Called(ctx, domain)

	return args.Error(0)
}

func (k *FakeKubernetes) IsCertificateProvisioned(ctx context.Context, secretName string) (bool, error) {
	args := k.Called(ctx, secretName)

	return args.Get(0).(bool), args.Error(1)
}
