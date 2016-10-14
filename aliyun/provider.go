package aliyun

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema:        providerSchema(),
		ResourcesMap:  providerResources(),
		ConfigureFunc: providerConfigure,
	}
}

// List of supported configuration fields for your provider.
// Here we define a linked list of all the fields that we want to
// support in our provider (api_key, endpoint, timeout & max_retries).
// More info in https://github.com/hashicorp/terraform/blob/v0.6.6/helper/schema/schema.go#L29-L142
func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_key": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "API Key used to authenticate with the service provider",
		},
		"endpoint": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "The URL to the API",
		},
		"timeout": &schema.Schema{
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Max. wait time we should wait for a successful connection to the API",
		},
		"max_retries": &schema.Schema{
			Type:        schema.TypeInt,
			Required:    true,
			Description: "The max. amount of times we will retry to connect to the API",
		},
	}
}

// List of supported resources and their configuration fields.
// Here we define da linked list of all the resources that we want to
// support in our provider. As an example, if you were to write an AWS provider
// which supported resources like ec2 instances, elastic balancers and things of that sort
// then this would be the place to declare them.
// More info here https://github.com/hashicorp/terraform/blob/v0.6.6/helper/schema/resource.go#L17-L81
func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"aliyun_instance": resourceAliyunInstance(),
	}
}

// This is the function used to fetch the configuration params given
// to our provider which we will use to initialise a dummy client that
// interacts with the API.
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := Config{
		APIKey:     d.Get("api_key").(string),
		Endpoint:   d.Get("endpoint").(string),
		Timeout:    d.Get("timeout").(int),
		MaxRetries: d.Get("max_retries").(int),
	}

	// You could have some field validations here, like checking that
	// the API Key is has not expired or that the username/password
	// combination is valid, etc.

	return &client, nil
}
