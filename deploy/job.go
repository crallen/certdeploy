package deploy

import (
	"io/ioutil"
	"net/http"

	"github.com/crallen/certdeploy/kubernetes"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type clusterJob struct {
	*clusterConfig
	Secrets    []*secretConfig
	kubeClient *kubernetes.Client
}

func (j *clusterJob) Run() {
	for _, s := range j.Secrets {
		for _, ns := range s.Namespaces {
			logger := log.WithFields(log.Fields{
				"cluster":   j.Name,
				"namespace": ns,
			})
			secret, err := j.kubeClient.Secret(s.Name, ns)
			if err == nil {
				if !j.setData(secret, s.Files, logger) {
					continue
				}
				_, err := j.kubeClient.UpdateSecret(secret, ns)
				if err != nil {
					logger.Error(err)
					continue
				}
				logger.Infof("updated secret %s", secret.Name)
			} else if statusErr, ok := err.(*errors.StatusError); ok && statusErr.Status().Code == http.StatusNotFound {
				secret := &v1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name: s.Name,
					},
					Type: v1.SecretTypeTLS,
				}
				if !j.setData(secret, s.Files, logger) {
					continue
				}
				_, err := j.kubeClient.CreateSecret(secret, ns)
				if err != nil {
					logger.Error(err)
					continue
				}
				logger.Infof("created secret %s", secret.Name)
			} else {
				logger.Error(err)
			}
		}
	}
}

func (j *clusterJob) setData(secret *v1.Secret, files map[string]string, logger *log.Entry) bool {
	dataMap := make(map[string][]byte)
	hasErrors := false
	for k, v := range files {
		data, err := ioutil.ReadFile(v)
		if err != nil {
			logger.Error(err)
			hasErrors = true
			continue
		}
		dataMap[k] = data
	}
	if !hasErrors {
		secret.Data = dataMap
		return true
	}
	return false
}
