package ucs

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/digitalonus/terraform-provider-ucs/ucsclient"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func sharedConfigForRegion(region string) (interface{}, error) {
	if os.Getenv("UCS_IP_ADDRESS") == "" {
		return nil, fmt.Errorf("empty UCS_IP_ADDRESS")
	}
	if os.Getenv("UCS_USERNAME") == "" {
		return nil, fmt.Errorf("empty UCS_USERNAME")
	}
	if os.Getenv("UCS_PASSWORD") == "" {
		return nil, fmt.Errorf("empty UCS_PASSWORD")
	}

	tls := false
	if os.Getenv("UCS_TSLINSECURESKIPVERIFY") == "true" {
		tls = true
	}

	logLevel, _ := strconv.Atoi(os.Getenv("UCS_IP_ADDRESS"))

	config := ucsclient.Config{
		IpAddress:             os.Getenv("UCS_IP_ADDRESS"),
		Username:              os.Getenv("UCS_USERNAME"),
		Password:              os.Getenv("UCS_PASSWORD"),
		TslInsecureSkipVerify: tls,
		LogLevel:              logLevel,
		LogFilename:           os.Getenv("UCS_LOG_FILENAME"),
	}

	return config.Client(), nil
}
