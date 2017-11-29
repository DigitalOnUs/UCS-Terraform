package ucs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func init() {
	resource.AddTestSweepers("ucs_service_profile", &resource.Sweeper{
		Name: "ucs_service_profile",
		F:    testSweepProfiles,
	})

}

func testSweepProfiles(region string) error {
	return nil
}

func TestAccUCSProfile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPrecheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: testAccUCSCProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUCSProfileConfig(),
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(
					"ucs_server_profile.test_server", "name", "test_server"),
				),
			},
		},
	})
}

func testAccCheckUCSProfileConfig() string {
	return fmt.Sprintf(`
		resource "ucs_service_profile" "test_server" {
  	name                     = "test_server"
  	target_org               = "root-org"
  	service_profile_template = "some-service-profile-template"
  	vnic {
    	name  = "eth0"
    	cidr = "1.2.3.4/24"
  	}
	}`)
}
