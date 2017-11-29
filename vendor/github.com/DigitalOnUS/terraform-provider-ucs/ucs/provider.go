package ucs

import (
	"github.com/DigitalOnUs/terraform-provider-ucs/ucsclient"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"ip_address": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UCS_IP_ADDRESS", nil),
				Description: "UCS Manager IP address or CIMC IP address.",
			},

			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UCS_USERNAME", nil),
				Description: "The user's name to access the UCS Management.",
			},

			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UCS_PASSWORD", nil),
				Description: "The password to access the UCS Management.",
			},

			"tslinsecureskipverify": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				DefaultFunc: schema.EnvDefaultFunc("UCS_TSLINSECURESKIPVERIFY", nil),
				Description: "The TSL insecure skip verify",
			},

			"log_level": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				DefaultFunc: schema.EnvDefaultFunc("UCS_LOG_LEVEL", nil),
				Description: "The log level",
			},

			"log_filename": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				DefaultFunc: schema.EnvDefaultFunc("UCS_LOG_FILENAME", nil),
				Description: "The log filename",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"ucs_service_profile": resourceUcsServiceProfile(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := ucsclient.Config{
		AppName:               "ucs",
		IpAddress:             d.Get("ip_address").(string),
		Username:              d.Get("username").(string),
		Password:              d.Get("password").(string),
		TslInsecureSkipVerify: d.Get("tslinsecureskipverify").(bool),
		LogFilename:           d.Get("log_filename").(string),
		LogLevel:              d.Get("log_level").(int),
	}

	return config.Client(), nil
}
