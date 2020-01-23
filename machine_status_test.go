package main

import (
	"testing"

	"github.com/OpenNebula/one/src/oca/go/src/goca"
)

type MockConnection struct {
}

func (c *MockConnection) GetAllVms() *goca.VMPool {
	vm1 := goca.VM{}
	vm1.Name = "sae-nexus"
	vm2 := goca.VM{}
	vm2.Name = "sae-devops-salt-testing"
	vmPool := goca.VMPool{}
	vmPool.VMs = append(vmPool.VMs, vm1)
	vmPool.VMs = append(vmPool.VMs, vm2)
	return &vmPool
}
func (c *MockConnection) GetVMIP(ID uint) string {
	return "10.12.161.10"
}

func (c *MockConnection) DestroyVM(ID uint) error {
	return nil
}

func (c *MockConnection) StartVM(name string, rolesString string, templateNum uint, hold bool) (string, uint) {
	return "10.12.161.10", 499
}

func TestGetAllRunningMachines(t *testing.T) {
	c := MockConnection{}
	pool := GetAllRunningMachines(&c)
	if len(pool.VMs) != 2 {
		t.Error("Not enough VMS found")
	}
}

func TestSetMachineStatus(t *testing.T) {
	//TODO mock OPEN_NEBULA_UP_ENV interface
	env := Env{}
	env.environment = ""
	c := MockConnection{}
	var mapList []map[string]string
	mach1 := Machine{}
	machs := make(map[string]OpenNebulaMachine)
	role1 := map[string]string{"name": "role-sae-rtapp", "tag_group_name": "SAE Systems"}
	mapList = append(mapList, role1)
	tag := Tags{mapList}
	machs["sae-nexus"] = mach1.GetNewOpenNebulaMachine("sae-nexus", tag, "salt-xenial")
	machs["sae-limon-dev"] = mach1.GetNewOpenNebulaMachine("sae-nexus", tag, "salt-xenial")

	machs = SetMachinesStatus(&c, machs, env)
	if machs["sae-nexus"].Status != false {
		t.Error("Incorrect Machines status")
	}
	if machs["sae-devops-salt-testing"].Status != false {
		t.Error("Incorrect Machines status")
	}
}
