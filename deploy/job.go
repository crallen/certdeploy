package deploy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/crallen/certdeploy/kubernetes"
)

type clusterJob struct {
	*clusterConfig
	kubeClient *kubernetes.Client
}

func (j *clusterJob) Run() {
	for _, s := range j.Secrets {
		for _, ns := range s.Namespaces {
			secret, err := j.kubeClient.Secret(s.Name, ns)
			if err == nil {
				_, err := j.kubeClient.UpdateSecret(secret, ns)
				if err != nil {
					j.logError(err, ns)
				}
			} else if statusErr, ok := err.(*errors.StatusError); ok && statusErr.Status().Code == http.StatusNotFound {
				secret := &v1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name: s.Name,
					},
				}
				for _, f := range s.Files {
					data, err := ioutil.ReadFile(f.Filename)
					if err != nil {
						j.logError(err, ns)
						continue
					}
					secret.Data[f.Key] = data
				}
				_, err := j.kubeClient.CreateSecret(secret, ns)
				if err != nil {
					j.logError(err, ns)
				}
			} else {
				j.logError(err, ns)
			}
		}
	}
}

func (j *clusterJob) logError(err error, namespace string) {
	if len(namespace) > 0 {
		fmt.Fprint(os.Stderr, fmt.Errorf("[%s] (%s) %v\n", j.Name, namespace, err))
	} else {
		fmt.Fprint(os.Stderr, fmt.Errorf("[%s] %v\n", j.Name, err))
	}
}
