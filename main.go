package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/mixslice/terraform-provider-aliyun/aliyun"
)

func main() {
	opts := plugin.ServeOpts{
		ProviderFunc: aliyun.Provider,
	}
	plugin.Serve(&opts)
}
