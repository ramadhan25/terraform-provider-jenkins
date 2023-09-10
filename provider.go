package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"jenkins_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jenkins_username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jenkins_api_token": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"jenkins_role": resourceJenkinsRole(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	url := d.Get("jenkins_url").(string)
	username := d.Get("jenkins_username").(string)
	apiToken := d.Get("jenkins_api_token").(string)

	config := MyProviderConfig{
		jenkins_url:       url,
		jenkins_username:  username,
		jenkins_api_token: apiToken,
	}
	return config, nil
}

type MyProviderConfig struct {
	jenkins_url       string
	jenkins_username  string
	jenkins_api_token string
}
