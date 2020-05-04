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
		Clusters: []*clusterConfig{
			{
				Name:    "dev",
				Context: "dev-cluster",
				Secrets: []*secretConfig{
					{
						Name: "tls-dev",
						Files: []*secretFile{
							{
								Key:      "fullchain.pem",
								Filename: "/etc/letsencrypt/live/myawesomedevdomain.com/fullchain.pem",
							},
							{
								Key:      "privkey.pem",
								Filename: "/etc/letsencrypt/live/myawesomedevdomain.com/privkey.pem",
							},
						},
						Namespaces: []string{
							"kube-system",
							"my-ns",
						},
					},
				},
			},
		},
	}
	assert.EqualValues(t, expected, cfg)
}
