package deploy

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
