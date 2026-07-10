package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/toby1991/go-zero/core/logx"
	"github.com/toby1991/go-zero/internal/devserver"
)

func TestServiceConf(t *testing.T) {
	c := ServiceConf{
		Name: "foo",
		Log: logx.LogConf{
			Mode: "console",
		},
		Mode: "dev",
		DevServer: devserver.Config{
			Port:       6470,
			HealthPath: "/healthz",
		},
	}
	c.MustSetUp()
}

func TestServiceConfWithMetricsUrl(t *testing.T) {
	c := ServiceConf{
		Name: "foo",
		Log: logx.LogConf{
			Mode: "volume",
		},
		Mode:       "dev",
		MetricsUrl: "http://localhost:8080",
	}
	assert.NoError(t, c.SetUp())
}
