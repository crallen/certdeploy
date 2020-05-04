package deploy

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/crallen/certdeploy/kubernetes"
)

type Runner struct {
	deployConfig *deployConfig
	kubeConfig   string
	jobC         chan *clusterJob
	errorC       chan error
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
			fmt.Fprintf(os.Stderr, "[%s] %v\n", c.Name, err)
			continue
		}
		jobC <- &clusterJob{
			clusterConfig: c,
			kubeClient:    kubeClient,
		}
		wg.Add(1)
	}

	wg.Wait()
	cancel()

	return nil
}
