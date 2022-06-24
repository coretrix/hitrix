# Kubernetes

This service is used for calling Kubernetes APIs.  
Currently, this service can be used to add other domains to Kubernetes Ingresses and handling
certification generation with cert-manager.

To use this, register the service into your `main.go` file first:

```go
registry.ServiceProviderKubernetes()
```

and you should put your credentials and other configs in your config file

```yml
kubernetes:
  environment: "dev"
  config_file: "/path/to/kubeconfig"
```
the `config_file` can be one of these:
- absolute path to config file
- relative path to config file - then address of config directory will be prepended to it 
- can be omitted completely - then Kubernetes in-cluster config of [Service Account Token](https://kubernetes.io/docs/admin/authentication/#service-account-tokens) will be used

Access the service:

```go
service.DI().Kubernetes()
```

Some functions this service provide are:

```go
	GetIngressDomains(ctx context.Context) ([]string, error)
	AddIngress(ctx context.Context, domain, secretName, serviceName, servicePortName string, annotations map[string]string) error
   	RemoveIngress(ctx context.Context, domain string) error
   	IsCertificateProvisioned(ctx context.Context, secretName string) (bool, error)
```
