package deploy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := loadConfig("../test/valid-config.yml")
	if err != nil {
		t.Error(err)
	}
	expected := &deployConfig{
		Secrets: map[string]*secretConfig{
			"tls-dev": {
				Name: "tls-dev",
				Files: map[string]string{
					"tls.crt": "/etc/letsencrypt/live/myawesomedevdomain.com/fullchain.pem",
					"tls.key": "/etc/letsencrypt/live/myawesomedevdomain.com/privkey.pem",
				},
				Namespaces: []string{
					"kube-system",
					"my-ns",
				},
			},
		},
		Clusters: []*clusterConfig{
			{
				Name:    "dev",
				Context: "dev-cluster",
				Secrets: []string{
					"tls-dev",
				},
			},
		},
	}
	assert.EqualValues(t, expected, cfg)
}
