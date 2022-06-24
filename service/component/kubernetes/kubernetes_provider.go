package kubernetes

import "context"

type IProvider interface {
	GetIngressDomains(ctx context.Context) ([]string, error)
	AddIngress(ctx context.Context, domain, secretName, serviceName, servicePortName string, annotations map[string]string) error
	RemoveIngress(ctx context.Context, domain string) error
	IsCertificateProvisioned(ctx context.Context, secretName string) (bool, error)
}
