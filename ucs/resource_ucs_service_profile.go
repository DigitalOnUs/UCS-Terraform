package ucs

import (
	"fmt"
	"net"
	"sync"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/thetonymaster/ucsclient"
	"github.com/thetonymaster/ucsclient/ipman"
)

type sessionCallback func(*ucsclient.UCSClient) error

var sessionMutex = sync.Mutex{}

func resourceUcsServiceProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceUcsServiceProfileCreate,
		Read:   resourceUcsServiceProfileRead,
		Update: resourceUcsServiceProfileUpdate,
		Delete: resourceUcsServiceProfileDelete,
		Schema: map[string]*schema.Schema{
			"service_profile_template": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"target_org": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"dn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Freestyle metadata for your resource",
			},
			"vnic": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"cidr": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"mac": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Creates a new Service Profile using the information available in the Resource Data.
// `meta` in this case is a pointer to a ucsclient.UCSClient.
func resourceUcsServiceProfileCreate(d *schema.ResourceData, meta interface{}) error {
	sp := &ucsclient.ServiceProfile{
		Name:         d.Get("name").(string),
		Template:     d.Get("service_profile_template").(string),
		TargetOrg:    d.Get("target_org").(string),
		VNICs:        make([]ucsclient.VNIC, 0, 1),
		Hierarchical: false,
	}

	vnics := d.Get("vnic").(*schema.Set).List()
	fmt.Printf("[INFO] vnics %+v\n", vnics)
	for _, item := range vnics {
		vnic := item.(map[string]interface{})
		sp.VNICs = append(sp.VNICs, ucsclient.VNIC{
			Name: vnic["name"].(string),
			CIDR: vnic["cidr"].(string),
		})

		// Validate the vnic's CIDR and return error if anything.
		err := validateCIDR(sp.VNICs[len(sp.VNICs)-1].CIDR)
		if err != nil {
			return err
		}
	}

	c := meta.(*ucsclient.UCSClient)

	err := withSession(c, func(client *ucsclient.UCSClient) error {
		d.Partial(true)
		if d.HasChange("name") {
			fmt.Printf("[INFO] Creating Profile \"%s\" from template \"%s\"\n", sp.Name, sp.Template)
			err := client.CreateServiceProfile(sp)
			if err != nil {
				fmt.Printf("[WARN] Failed to create profile \"%s\": %s\n", sp.Name, err)
				return err
			}

			fmt.Printf("[INFO] Profile \"%s\" was created\n", sp.Name)
			d.SetId(sp.Name) // tell Terraform that a profile was created. The existence of a non-blank ID is what tells Terraform that a profile was created
			d.Set("dn", sp.DN())
			d.SetPartial("name")
		}

		if d.HasChange("vnic") {
			vnics := make([]map[string]string, len(sp.VNICs))
			// Assign an IP to each of the vnics in the Service Profile.
			for i, vnic := range sp.VNICs {
				ip, err := ipman.GenerateIP("inventory-"+vnic.Name, vnic.CIDR)
				if err != nil {
					return err
				}
				vnic.Ip = ip

				vnics[i] = map[string]string{
					"name": vnic.Name,
					"ip":   vnic.Ip.String(),
					"cidr": vnic.CIDR,
				}
			}
			d.Set("vnic", vnics)
			d.SetPartial("vnic")
		}

		d.Partial(false)
		fmt.Printf("[DEBUG] Exiting resourceUcsServiceProfileCreate(...)\n")
		return nil
	})

	if err != nil {
		return err
	}

	return resourceUcsServiceProfileRead(d, c)
}

// Fetches general information of the Service Profile from UCS.
// If the Service Profile information has been changed remotely
// this method will update the local state file (terraform.tfstate).
// If the Service Profile is no longer available this will remove it from
// the tfstate file.
func resourceUcsServiceProfileRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ucsclient.UCSClient)

	err := withSession(c, func(client *ucsclient.UCSClient) error {
		fmt.Printf("[DEBUG] Entering resourceUcsServiceProfileRead(...)\n")

		//1. Query the UCS for the profile
		dn := d.Get("dn").(string)
		sp, err := client.ConfigResolveDN(dn)

		if err != nil {
			return err
		}

		// If the service profile could not be found we assume that it does not exist anymore
		// We tell Terraform so by setting its id to a blank string.
		if sp == nil {
			d.SetId("")
			return fmt.Errorf("Service Profile: %s", "No longer exists")
		}

		// Fetch vnic info from ResourceData
		vNicsFromResourceData := fetchVnicsFromResourceData(d)

		// Merge the UCS vnic info with the ResourceData vnic info
		vnics := mergeVnics(vNicsFromResourceData, sp.VNICs)

		// Update the information related to the service profile fetched from UCS in Terraform.
		d.Set("name", sp.Name)
		d.Set("service_profile_template", sp.Template)
		d.Set("target_org", sp.TargetOrg)
		d.Set("vnic", vnics)

		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": d.Get("vnic.0.ip").(string),
		})

		fmt.Printf("[DEBUG] Exiting resourceUcsServiceProfileRead(...)\n")
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// Updates the Service Profile in UCS.
func resourceUcsServiceProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ucsclient.UCSClient)
	fmt.Printf("[DEBUG] Entering resourceUcsServiceProfileUpdate(...)\n")
	fmt.Printf("[DEBUG] Exiting resourceUcsServiceProfileUpdate(...)\n")
	return resourceUcsServiceProfileRead(d, c)
}

// Deletes a given Service Profile, using its "dn" as the identifier.
func resourceUcsServiceProfileDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ucsclient.UCSClient)
	return withSession(c, func(client *ucsclient.UCSClient) error {
		fmt.Printf("[DEBUG] Entering resourceUcsServiceProfileDelete(...)\n")

		name := d.Id()
		targetOrg := d.Get("target_org").(string)

		// Delete the resource
		err := client.Destroy(name, targetOrg, true)
		if err != nil {
			return err
		}

		// Tell Terraform that the resource has been successfully destroyed
		d.SetId("")

		fmt.Printf("[DEBUG] Exiting resourceUcsServiceProfileDelete(...)\n")
		return nil
	})
}

func fetchVnicsFromResourceData(d *schema.ResourceData) (ret []ucsclient.VNIC) {
	vnics := d.Get("vnic").(*schema.Set).List()
	for _, item := range vnics {
		vnic := item.(map[string]interface{})
		ret = append(ret, ucsclient.VNIC{
			Name: vnic["name"].(string),
			Ip:   net.ParseIP(vnic["ip"].(string)),
			CIDR: vnic["cidr"].(string),
		})
	}
	return
}

func mergeVnics(localVnics, remoteVnics []ucsclient.VNIC) (ret []map[string]string) {
	for _, local := range localVnics {
		for _, remote := range remoteVnics {
			if local.Name == remote.Name {
				ret = append(ret, map[string]string{
					"name": remote.Name,
					"ip":   local.Ip.String(),
					"mac":  remote.Mac,
					"cidr": local.CIDR,
				})
				break
			}
		}
	}
	return
}

func validateCIDR(cidr string) (err error) {
	_, _, err = net.ParseCIDR(cidr)
	return
}

func withSession(c *ucsclient.UCSClient, cb sessionCallback) error {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	err := c.Login()
	if err != nil {
		return err
	}

	cb(c)
	c.Logout()
	return nil
}
