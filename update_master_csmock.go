package main

import (
	"fmt"
	"os"
	"os/exec"
)

func UpdateCsmock(env Env) error {
	path := env.config.ConfigPath + "OpenNebulaUpCsmock.yml"
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, mach := range env.machs {
		mach.AddToCsmock(f)
	}
	PushCsmockToMaster(env.machs["salt-master"].IP, env)
	PushBaseCsmockToMaster(env.machs["salt-master"].IP, env)
	return nil
}

func (mach OpenNebulaMachine) AddToCsmock(f *os.File) error {
	if mach.Status {
		f.WriteString(fmt.Sprintln(mach.Hostname + ":"))
		f.WriteString(fmt.Sprintf("  name: %s\n", mach.Hostname))
		f.WriteString(fmt.Sprintln("  tags:"))
		for _, data := range mach.TagData.Tags {
			f.WriteString(fmt.Sprintf("    - name: %s\n", data["name"]))
			f.WriteString(fmt.Sprintf("      tag_group_name: %s\n", data["tag_group_name"]))
		}
		f.WriteString("\n")
	}
	return nil
}

func PushBaseCsmockToMaster(IP string, env Env) error {
	dest := fmt.Sprintf("root@%s:/srv/states/production/vagrant-env/", IP)
	path := env.config.ProjectsPath + "/salt-states/vagrant-env/"
	cmd := exec.Command("scp", "-r", path, dest)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error in scp of base csmock")
		panic(err)
	}
	return nil
}

func PushCsmockToMaster(IP string, env Env) error {
	dest := fmt.Sprintf("root@%s:/srv/states/production/csmock.yml", IP)

	path := env.config.ConfigPath + "OpenNebulaUpCsmock.yml"
	cmd := exec.Command("scp", path, dest)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error in scp of OpenNebulUpCsmock csmock")
		panic(err)
	}

	return nil
}
