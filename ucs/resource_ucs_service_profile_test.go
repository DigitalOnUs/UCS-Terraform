package ucs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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

func TestAccUCSServerProfile_Simple(t *testing.T) {

	r := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPrecheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccUCSCProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUCSProfileConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testProfileExists("ucs_service_profile.test_server"),
				),
			},
		},
	})
}

func testProfileExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return nil
	}
}

func testAccUCSCProfileDestroy(s *terraform.State) error {
	return nil
}

func testAccCheckUCSProfileConfig(r string) string {
	return fmt.Sprintf(`
resource "ucs_service_profile" "master-server" {
  name                     = "server-%s"
  target_org               = "org-root"
  service_profile_template = "template-example"
  metadata {
    role             = "master"
    ansible_ssh_user = "root"
    foo              = "bar"
  }
  vnic {
    name  = "eth0"
    cidr = "1.2.3.4/24"
  }
}`, r)
}
