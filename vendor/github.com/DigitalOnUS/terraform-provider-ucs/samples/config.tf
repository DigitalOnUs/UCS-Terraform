provider "ucs" {
  ip_address   = "172.16.63.137"
  username     = "ucspe"
  password     = "ucspe"
  log_level    = 1
  log_filename = "terraform.log"
  tslinsecureskipverify = true
}

resource "ucs_service_profile" "master-server" {
  name                     = "Server 3"
  target_org               = "root-org"
  service_profile_template = "terraformprofiletemplate"
  metadata { # This field is pretty much free style. Values must always be strings.
    role             = "master" # This is useful when creating a Mantl cluster
    ansible_ssh_user = "root"
    foo              = "bar"
  }
  vnic {
  name  = "eth0"
  cidr = "1.2.3.4/24"
}
}
