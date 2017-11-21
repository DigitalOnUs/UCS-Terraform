package ucs

import (
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
					"cloudscale_server.basic", "flavor_slug", "flex-2"),
				),
			},
		},
	})
}

func testAccCheckUCSProfileConfig() string {
	return `
	resource "ucs_service_profile" "the-server-name" {
  	name                     = "the-server-name"
  	target_org               = "some-target-org"
  	service_profile_template = "some-service-profile-template"
  	vNIC {
    	name  = "eth0"
    	cidr = "1.2.3.4/24"
  	}
	}`
}
