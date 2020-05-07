package deploy

import (
	"context"
	"sync"

	"github.com/crallen/certdeploy/kubernetes"
	log "github.com/sirupsen/logrus"
)

type Runner struct {
	deployConfig *deployConfig
	kubeConfig   string
}

func New(configFile, kubeConfig string) (*Runner, error) {
	config, err := loadConfig(configFile)
	if err != nil {
		return nil, err
	}
	return &Runner{
		deployConfig: config,
		kubeConfig:   kubeConfig,
	}, nil
}

func (r *Runner) Run() error {
	var wg sync.WaitGroup

	jobC := make(chan *clusterJob, len(r.deployConfig.Clusters))
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		for {
			select {
			case job := <-jobC:
				job.Run()
				wg.Done()

			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	for _, c := range r.deployConfig.Clusters {
		kubeClient, err := kubernetes.New(c.Context, r.kubeConfig)
		if err != nil {
			log.WithField("cluster", c.Name).Error(err)
			continue
		}
		secrets := make([]*secretConfig, 0)
		for _, s := range c.Secrets {
			if secret, ok := r.deployConfig.Secrets[s]; ok {
				secrets = append(secrets, secret)
			} else {
				log.WithField("cluster", c.Name).Errorf("no configuration found for secret %s", s)
			}
		}
		wg.Add(1)
		jobC <- &clusterJob{
			clusterConfig: c,
			Secrets:       secrets,
			kubeClient:    kubeClient,
		}
	}

	wg.Wait()
	cancel()

	return nil
}
