package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type ioData interface {
	ReadDir(dirname string) ([]string, error)
	WriteString(line string) error
	SetDataLocation(location string) error
	Create() error
}

type fileIO struct {
	file string
}

func (data fileIO) Create() error {
	_, err := os.Create(data.file)
	if err != nil {
		fmt.Printf("Can't create file\n")
		panic(err)
	}
	return err
}

func (data fileIO) WriteString(line string) error {
	f, err := os.OpenFile(data.file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		fmt.Printf("Couldn't Open File: %s for writing\n", data.file)
	}
	_, err = f.WriteString(line)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (data *fileIO) SetDataLocation(location string) error {
	data.file = location
	return nil
}

func (data fileIO) ReadDir(dirname string) ([]string, error) {
	files, err := ioutil.ReadDir(dirname)
	var fileNames []string
	if err != nil {
		return fileNames, err
	}
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return fileNames, nil
}

func GetReposMachines(repo string) map[string]Machine {
	data := GetMachinesYaml(repo)
	return ReadYamlData(data)
}
func GetReposTags(repo string) map[string]Tags {
	data := GetCsMock(repo)
	return ReadCsMockData(data)
}

func GetMachinesYaml(repo string) []byte {
	path := GetMachinesYamlPath(repo)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte("")
	}
	return data
}

func GetCsMock(repo string) []byte {
	path := GetCsmockPath(repo)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte("")
	}
	return data
}

func ReadYamlData(data []byte) map[string]Machine {
	machines := make(map[string]Machine)
	err := yaml.Unmarshal(data, &machines)
	if err != nil {
		fmt.Printf("Error unmarshalling machins.yml yaml in %s\n", string(data))
		panic(err)
	}
	return machines
}

func ReadCsMockData(data []byte) map[string]Tags {
	tags := make(map[string]Tags)
	err := yaml.Unmarshal(data, &tags)
	if err != nil {
		fmt.Printf("Error unmarshalling csmock yaml in %s\n", string(data))
		panic(err)
	}
	return tags
}

func GetMachinesYamlPath(repo string) string {
	if repo == "salt-states" {
		return "../salt-states/vagrant-env/machines.yml"
	}
	base_machines := "../%s/machines.yml"
	return fmt.Sprintf(base_machines, repo)
}

func GetCsmockPath(repo string) string {
	if repo == "salt-states" {
		return ""
	}
	base_machines := "../%s/csmock.yml"
	return fmt.Sprintf(base_machines, repo)
}

func FindAllSaltDirectories(data ioData) []string {
	files, err := data.ReadDir("../")
	if err != nil {
		log.Fatal(err)
	}
	var dirNames []string
	for _, file := range files {
		if strings.HasPrefix(file, "salt-") && file != "salt-states" {
			dirNames = append(dirNames, file)
		}
	}
	return dirNames
}

func GetAllMachines(data ioData) map[string]OpenNebulaMachine {
	machs := make(map[string]OpenNebulaMachine)
	repos := FindAllSaltDirectories(data)
	var machines map[string]Machine
	var tags map[string]Tags
	for _, repo := range repos {
		if repo == "salt-pillar" {
			continue // no machines in pillar
		}
		machines = GetReposMachines(repo)
		tags = GetReposTags(repo)
		for name, machine := range machines {
			tag := GetTagLookup(tags, name) // for partial matching of tags
			operatingSystem := GetOperatingSystem(machines, name)
			machs[name] = machine.GetNewOpenNebulaMachine(name, tag, operatingSystem)
		}
	}
	base_machines := GetReposMachines("salt-states")
	for name, machine := range base_machines {
		var empty_tag = Tags{}
		operatingSystem := GetOperatingSystem(base_machines, name)
		machs[name] = machine.GetNewOpenNebulaMachine(name, empty_tag, operatingSystem)
	}

	return machs
}

func GetTagLookup(tags map[string]Tags, name string) Tags {
	var machine_tag = Tags{}
	if tag, ok := tags[name]; ok {
		for _, tag_data := range tag.Tags {
			machine_tag.Tags = append(machine_tag.Tags, tag_data)
		}
	}
	for tag_title, tag_object := range tags {
		if strings.HasSuffix(tag_title, "-") && strings.HasPrefix(name, tag_title) {
			for _, tag_data := range tag_object.Tags {
				machine_tag.Tags = append(machine_tag.Tags, tag_data)
			}
		}
	}
	return machine_tag

}
func (mach *Machine) GetNewOpenNebulaMachine(name string, tag_data Tags, operatingSystem string) OpenNebulaMachine {
	var roles []string
	for _, data := range tag_data.Tags {
		if strings.HasPrefix(data["name"], "role-") {
			roles = append(roles, strings.TrimLeft(data["name"], "role-"))
		}
	}
	var team_env string
	var MySyncedFolders []SyncedFolder
	if mach.TeamEnv == "" {
		team_env = "qa"
	} else {
		team_env = mach.TeamEnv
	}
	for _, folder := range mach.SyncedFolder {
		SyncedFolder := SyncedFolder{folder["source"], folder["destination"]}
		MySyncedFolders = append(MySyncedFolders, SyncedFolder)
	}

	return OpenNebulaMachine{name, operatingSystem, roles, "", 0, false, mach.Hostname, mach.Mem, mach.VCPU, MySyncedFolders, team_env, tag_data}
}

func GetOperatingSystem(machines map[string]Machine, name string) string {
	inherits := machines[name].Inherit
	if inherits == "" {
		return name
	}
	if _, ok := machines[inherits]; !ok {
		return inherits
	}
	return GetOperatingSystem(machines, inherits)
}
