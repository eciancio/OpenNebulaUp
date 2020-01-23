package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/OpenNebula/one/src/oca/go/src/goca"
)

var machineMap = map[string]uint{"salt-xenial": 116, "salt-bionic": 119, "salt-trusty": 118, "salt-precise": 120, "salt-centos": 121, "salt-freebsd": 122}

func (mach OpenNebulaMachine) StartOpenNebulaMachine(env Env, saltMasterIP string, hold bool) OpenNebulaMachine {
	if mach.Status {
		fmt.Printf("%s is already up and running at %s\n", mach.Name, mach.IP)
		os.Exit(0)
	}

	template_number := mach.GetTemplateNumber()
	roles := mach.GetRolesString(saltMasterIP, env)
	mach.Status = true
	env.machs[mach.Name] = mach

	if env.machs["salt-master"].Status && mach.Name != "salt-master" {
		// updating csmock

		ClearKeys(env)
		UpdateCsmock(env)
	}
	mach.IP, mach.ID = env.connection.StartVM(mach.Name, roles, template_number, hold)
	mach.Status = true

	return mach
}

func (mach OpenNebulaMachine) PowerOffOpenNebulaMachine(env Env) error {
	if mach.Name == "salt-master" {
		StopModd(env)
	}
	fmt.Println("Powering Off " + mach.Name)
	return env.connection.PowerOffVM(mach.ID)
}
func (mach OpenNebulaMachine) DestroyOpenNebulaMachine(env Env) error {
	if mach.Name == "salt-master" {
		StopModd(env)
	}
	fmt.Println("Destroying " + mach.Name)
	return env.connection.DestroyVM(mach.ID)
}

func (mach OpenNebulaMachine) ResumeOpenNebulaMachine(env Env) error {
	if mach.Name == "salt-master" {
		StartModd(env)
	}
	fmt.Println("Resuming " + mach.Name)
	return env.connection.Resume(mach.ID)
}

func (mach OpenNebulaMachine) RebootOpenNebulaMachine(env Env) error {
	fmt.Println("Rebooting " + mach.Name)
	return env.connection.Resume(mach.ID)
}

func (mach OpenNebulaMachine) SuspendOpenNebulaMachine(env Env) error {
	return env.connection.SuspendVM(mach.ID)
}

func (con *gocaConnection) StartVM(name string, rolesString string, templateNum uint, hold bool) (string, uint) {
	goca.SetClient(con.config)
	template := goca.NewTemplate(templateNum)
	id, err := template.Instantiate(name, hold, rolesString, false)
	if err != nil {
		fmt.Printf("Error in Template initation of %s with roles %s\n", name, rolesString)
		panic(err)
	}

	minion := goca.NewVM(id)
	minion.Info()
	return minion.Template.NICs[0].IP, id
}

func (con *gocaConnection) Resume(ID uint) error {
	vm := goca.NewVM(ID)
	return vm.Resume()
}

func (con *gocaConnection) Reboot(ID uint) error {
	vm := goca.NewVM(ID)
	return vm.Reboot()
}

func (con *gocaConnection) DestroyVM(ID uint) error {
	vm := goca.NewVM(ID)
	return vm.TerminateHard()
}

func (con *gocaConnection) SuspendVM(ID uint) error {
	vm := goca.NewVM(ID)
	return vm.Suspend()
}

func (con *gocaConnection) PowerOffVM(ID uint) error {
	vm := goca.NewVM(ID)
	return vm.Poweroff()
}

func (mach OpenNebulaMachine) GetTemplateNumber() uint {
	if mach.Name == "salt-master" {
		return uint(117)
	}
	if strings.Contains(mach.OperatingSystem, "centos") {
		return machineMap["salt-centos"]
	}
	if strings.Contains(mach.OperatingSystem, "llbsd") || strings.Contains(mach.OperatingSystem, "epos") {
		return machineMap["salt-freebsd"]
	}
	if os, ok := machineMap[mach.OperatingSystem]; ok {
		return os
	}

	return uint(116)
}
func (mach *OpenNebulaMachine) GetRolesString(saltMasterIP string, env Env) string {
	publickey := "\"" + GetPublicKey(env) + "\""
	if mach.Name == "salt-master" {
		return "NODNS=YES opennebulaupenv=" + env.environment + " sshpublickey=" + publickey
	}
	roles := "roles=none"
	roles += " opennebulaupenv=" + env.environment
	roles += " saltmaster=" + saltMasterIP
	roles += " sshpublickey=" + publickey
	hostname := mach.GetHostname(env)
	hostname = strings.TrimSuffix(hostname, ".devnet.llnw.net")
	if hostname != "" {
		roles += " HOSTNAME=" + hostname
	} else {
		hostname = mach.Name
	}
	mach.Hostname = hostname + ".devnet.llnw.net"
	if mach.Mem != 0 {
		roles += " MEMORY=" + strconv.Itoa(mach.Mem)
	}
	if mach.VCPU != 0 {
		roles += " VCPU=" + strconv.Itoa(mach.VCPU)
	}
	roles += " TEAM_ENV=" + mach.TeamEnv

	return roles
}
