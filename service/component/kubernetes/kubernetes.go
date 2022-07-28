package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"strings"

	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Kubernetes struct {
	client        *kubernetes.Clientset
	dynamicClient dynamic.Interface
	namespace     string
}

func NewKubernetes(configFilePath, namespace string) *Kubernetes {
	// load config from file or create the in-cluster config if path is empty
	config, err := clientcmd.BuildConfigFromFlags("", configFilePath)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return &Kubernetes{
		client:        clientset,
		dynamicClient: dynamicClient,
		namespace:     namespace,
	}
}

func (k *Kubernetes) IsCertificateProvisioned(ctx context.Context, secretName string) (bool, error) {
	deploymentRes := schema.GroupVersionResource{Group: "cert-manager.io", Version: "v1", Resource: "certificates"}

	result, getErr := k.dynamicClient.Resource(deploymentRes).Namespace(k.namespace).Get(ctx, secretName, metav1.GetOptions{})
	if getErr != nil {
		return false, fmt.Errorf("failed to get the latest version of Deployment: %v", getErr)
	}

	if statusObject, exists := result.Object["status"]; exists {
		if statusObject, ok := statusObject.(map[string]interface{}); ok {
			if conditionsObject, exists := statusObject["conditions"]; exists {
				if conditionsObject, ok := conditionsObject.([]interface{}); ok && len(conditionsObject) > 0 {
					if conditionsObject, ok := conditionsObject[0].(map[string]interface{}); ok {
						status := conditionsObject["status"].(string)
						message := conditionsObject["message"].(string)

						if status == "True" {
							return true, nil
						}

						return false, errors.New(message)
					}
				}
			}
		}
	}

	return false, errors.New("unknown error happened. try again")
}

func (k *Kubernetes) GetIngressDomains(ctx context.Context) ([]string, error) {
	ingresses, err := k.client.NetworkingV1().Ingresses(k.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	domains := make([]string, 0, len(ingresses.Items))

	for _, ingress := range ingresses.Items {
		for _, tls := range ingress.Spec.TLS {
			domains = append(domains, tls.Hosts...)
		}
	}

	return domains, nil
}

func (k *Kubernetes) AddIngress(ctx context.Context, domain, secretName, serviceName, servicePortName string, annotations map[string]string) error {
	domain = strings.TrimSpace(strings.ToLower(domain))
	// TODO: Validate host

	pathType := networkingv1.PathTypeImplementationSpecific

	ingress := &networkingv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: networkingv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        domain,
			Namespace:   k.namespace,
			Annotations: annotations,
		},
		Spec: networkingv1.IngressSpec{
			TLS: []networkingv1.IngressTLS{
				{
					Hosts:      []string{domain},
					SecretName: secretName,
				},
			},
			Rules: []networkingv1.IngressRule{
				{
					Host: domain,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: serviceName,
											Port: networkingv1.ServiceBackendPort{
												Name: servicePortName,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Status: networkingv1.IngressStatus{},
	}

	_, err := k.client.NetworkingV1().Ingresses(k.namespace).Create(ctx, ingress, metav1.CreateOptions{})

	if apierrors.IsAlreadyExists(err) {
		// ignore already exist error
		return nil

		// TODO: or we can
		// return fmt.Errorf("%s has already been registered", domain)
	}

	return err
}

func (k *Kubernetes) RemoveIngress(ctx context.Context, domain string) error {
	domain = strings.TrimSpace(strings.ToLower(domain))
	// TODO: Validate host

	err := k.client.NetworkingV1().Ingresses(k.namespace).Delete(ctx, domain, metav1.DeleteOptions{})

	if apierrors.IsNotFound(err) {
		// ignore already exist error
		return nil

		// TODO: or we can
		// return fmt.Errorf("%s has not been registered", domain)
	}

	return err
}
