package aliyun

import (
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/oss"
)

type Config struct {
	AccessKey string
	SecretKey string
	Region    common.Region
}

type AliyunClient struct {
	ecsclient *ecs.Client
	ossclient *oss.Client
	region    common.Region
}

// Client configures and returns a fully initialized Client
func (c *Config) Client() (interface{}, error) {
	var client AliyunClient

	client.ecsclient = ecs.NewClient(c.AccessKey, c.SecretKey)
	// client.ossclient = oss.NewOSSClient()
	client.region = c.Region

	return &client, nil
}
