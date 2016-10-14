package aliyun

import (
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/oss"
)

type Config struct {
	AccessKey string
	SecretKey string
}

type AliyunClient struct {
	ecsclient *ecs.Client
	ossclient *oss.Client
}

// Client configures and returns a fully initialized Client
func (c *Config) Client() (interface{}, error) {
	var client AliyunClient

	client.ecsclient = ecs.NewClient(c.AccessKey, c.SecretKey)
	// client.ossclient = oss.NewOSSClient()

	return &client, nil
}
