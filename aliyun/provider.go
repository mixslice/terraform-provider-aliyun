package aliyun

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["access_key"],
			},
			"secret_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["secret_key"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"aliyun_ecs_image": dataSourceAliyunEcsImage(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"aliyun_ecs_instance": resourceAliyunEcsInstance(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key": "The access key for API operations. You can retrieve this\n" +
			"from the 'AccessKeys' section of the Aliyun console.",

		"secret_key": "The secret key for API operations. You can retrieve this\n" +
			"from the 'AccessKeys' section of the Aliyun console.",
	}
}

// This is the function used to fetch the configuration params given
// to our provider which we will use to initialise a dummy client that
// interacts with the API.
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AccessKey: d.Get("access_key").(string),
		SecretKey: d.Get("secret_key").(string),
	}

	// You could have some field validations here, like checking that
	// the API Key is has not expired or that the username/password
	// combination is valid, etc.

	return config.Client()
}
