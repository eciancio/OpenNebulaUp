package main

import (
	"testing"
)

var hostname_yml = `
msop.phx2:
  roles:
  - role-sae-msops

nexus-01.phx3:
  roles:
  - role-sae-nexus
  - role-sae-devops-salt-testing
  - role-sae-singel
`

func TestReadHostnameData(t *testing.T) {
	names := ReadHostnameData([]byte(hostname_yml))
	if machine, ok := names["msop.phx2"]; ok {
		if len(machine.Roles) != 1 {
			t.Errorf("Number of roles is wrong")
		} else if machine.Roles[0] != "role-sae-msops" {
			t.Errorf("Role is wrong")
		}
	} else {
		t.Errorf("Machine not available")
	}
	if machine, ok := names["nexus-01.phx3"]; ok {
		if len(machine.Roles) != 3 {
			t.Errorf("Number of roles is wrong")
		} else if machine.Roles[0] != "role-sae-nexus" {
			t.Errorf("Role is wrong")
		}
	} else {
		t.Errorf("Machine not available")
	}
}

func TestGetRolesIndex(t *testing.T) {
	roles := []string{"role-sae-nexus", "role-sae-singel", "role-sae-devops-salt-testing"}
	index := GetRolesIndex(roles)
	if index != "role-sae-devops-salt-testing,role-sae-nexus,role-sae-singel" {
		t.Error("Did not index correctly")
	}
	role2 := []string{"role-sae-singel", "role-sae-devops-salt-testing", "role-sae-nexus"}
	index2 := GetRolesIndex(role2)
	if index != index2 {
		t.Error("Two Indexes that should match don't")
	}
}

func TestGetHostnameMapping(t *testing.T) {
	roles := []string{"role-sae-nexus", "role-sae-singel", "role-sae-devops-salt-testing"}
	index := GetRolesIndex(roles)
	hostNameMap := GetHostnameMap([]byte(hostname_yml))
	if hostname, ok := hostNameMap[index]; !ok {
		t.Error("Index Not Present")
	} else if hostname != "nexus-01.phx3" {
		t.Error("Hostname incorrect")
	}
}
