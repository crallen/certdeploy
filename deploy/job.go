package deploy

import "github.com/crallen/certdeploy/kubernetes"

type clusterJob struct {
	config     *clusterConfig
	kubeClient *kubernetes.Client
}
