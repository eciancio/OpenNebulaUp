package main

import (
	"testing"
)

var machines_yml = `
salt-master:
 inherit: salt-trusty
 enabled: True
 autostart: True
 mem: 2048
cat-master:
 inherit: salt-trusty
 enabled: True
 autostart: True
 mem: 2048
 hostname: mastering-cats
 vcpu: 23
 base_env: dev
 synced_folder:
  - {source: "/Users/.m2/repository", destination: "/tmp/.m2/repository"}
`

var csmock_yml = `
sae-msops:
  name: sae-msops
  building_name_short: YXX1
  tags:
    - name: role-sae-msops
      tag_group_name: vagrant

sae-pxemaster:
  name: sae-pxemaster
  tags:
    - name: role-sae-pxemaster
      tag_group_name: vagrant
    - name: role-sae-pxemaster3
      tag_group_name: vagrant
`

type MockIO struct {
	data string
}

func (r *MockIO) Create() error {
	return nil
}

func (r *MockIO) ReadDir(dirname string) ([]string, error) {
	var files []string
	files = append(files, "salt-states")
	files = append(files, "not-salt")
	files = append(files, "saltfake")
	files = append(files, "salt-sae")
	files = append(files, "salt-pillar")
	return files, nil
}

func (r *MockIO) SetDataLocation(location string) error {
	return nil
}

func (r *MockIO) WriteString(data string) error {
	r.data += data
	return nil
}

func TestReadYamlData(t *testing.T) {
	machines := ReadYamlData([]byte(machines_yml))
	if _, ok := machines["salt-master"]; !ok {
		t.Fail()
	}
	if machines["salt-master"].Inherit != "salt-trusty" {
		t.Fail()
	}
	if machines["salt-master"].Mem != 2048 {
		t.Fail()
	}
	if machines["cat-master"].VCPU != 23 {
		t.Error("Incorrect number of VCPU")
	}
	if machines["cat-master"].Hostname != "mastering-cats" {
		t.Error("Incorrect Hostname")
	}
	if machines["cat-master"].SyncedFolder[0]["destination"] != "/tmp/.m2/repository" {
		t.Error("Incorrect Sync Destination")
	}

}
func TestReadCsmockData(t *testing.T) {
	tags := ReadCsMockData([]byte(csmock_yml))
	if _, ok := tags["sae-msops"]; !ok {
		t.Error("sae-msops not read from csmock")
	}
	if tags["sae-msops"].Tags[0]["name"] != "role-sae-msops" {
		t.Error("Tag not found for sae-msops")
	}
	if len(tags["sae-pxemaster"].Tags) != 2 {
		t.Error("Two tags not read for sae-pxemaster")
	}
}

func TestGetYamlPath(t *testing.T) {
	path := GetMachinesYamlPath("salt-states")
	if path == "vagrant-env/machines.yml" {
		t.Error("salt-states path is incorrect")
	}
	path = GetMachinesYamlPath("salt-sae")
	if path != "../salt-sae/machines.yml" {
		t.Error("salt-sae vagrant path is incorrect")
	}
}

func TestGetCsmockPath(t *testing.T) {
	path := GetCsmockPath("salt-states")
	if path != "" {
		t.Error("salt-states path is incorrect")
	}
	path = GetCsmockPath("salt-sae")
	if path != "../salt-sae/csmock.yml" {
		t.Error("salt-sae csmock path is incorrect")
	}
}
func TestFindAllSaltDirectories(t *testing.T) {
	data := MockIO{}
	dirs := FindAllSaltDirectories(&data)
	if len(dirs) != 2 {
		t.Error("Incorrect number of directories counted")
	}
}

func TestGetNewOpenNebulaMachine(t *testing.T) {
	var mapList []map[string]string
	role1 := map[string]string{"name": "role-sae-rtapp", "tag_group_name": "SAE Systems"}
	role2 := map[string]string{"name": "role-sae-nexus", "tag_group_name": "SAE Systems"}
	mapList = append(mapList, role1)
	mapList = append(mapList, role2)
	tag := Tags{mapList}
	mach1 := Machine{}
	machine := mach1.GetNewOpenNebulaMachine("sae-nexus", tag, "salt-xenial")
	if machine.Name != "sae-nexus" {
		t.Error("Incorrect machine name")
	}
	if machine.OperatingSystem != "salt-xenial" {
		t.Error("Incorrect OpperatingSystem")
	}
	if len(machine.Roles) != 2 {
		t.Error("Incorrect Number of Roles")
	}
}

func TestGetOperatingSystem(t *testing.T) {
	machs := make(map[string]Machine)
	dummy := make([]map[string]string, 0)
	machs["sae-limon"] = Machine{"salt-xenial", true, true, 2048, "", 0, "", dummy}
	machs["sae-limon-dev"] = Machine{"sae-limon", true, true, 2048, "", 0, "", dummy}
	os := GetOperatingSystem(machs, "sae-limon-dev")
	if os != "salt-xenial" {
		t.Error("Incorrect OpperatingSystem")
	}

}
