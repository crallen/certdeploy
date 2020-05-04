package kubernetes

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	kubeClient *kubernetes.Clientset
}

func New(contextName, kubeConfig string) (*Client, error) {
	cfg, err := loadConfig(contextName, kubeConfig)
	if err != nil {
		return nil, err
	}
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{kubeClient}, nil
}

func (c *Client) Secret(name string, namespace string) (*v1.Secret, error) {
	return c.kubeClient.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
}

func (c *Client) CreateSecret(secret *v1.Secret, namespace string) (*v1.Secret, error) {
	return c.kubeClient.CoreV1().Secrets(namespace).Create(secret)
}

func (c *Client) UpdateSecret(secret *v1.Secret, namespace string) (*v1.Secret, error) {
	return c.kubeClient.CoreV1().Secrets(namespace).Update(secret)
}

func loadConfig(contextName, kubeConfig string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfig},
		&clientcmd.ConfigOverrides{
			CurrentContext: contextName,
		},
	).ClientConfig()
}
