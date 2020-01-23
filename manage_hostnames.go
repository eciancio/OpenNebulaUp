package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

func (mach OpenNebulaMachine) ReadHostname(env Env) string {
	hostname := mach.Hostname
	if hostname != "" {
		return hostname
	}

	hostname = mach.GetRepresentative(env)
	if hostname != "" {
		return hostname
	}
	return ""
}

func (mach OpenNebulaMachine) GetHostname(env Env) string {
	hostname := mach.ReadHostname(env)
	if hostname != "" {
		return hostname
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Representative Hostname not found for VM\nIf no specific hostname is needed for top.sls hit Enter\nElse type in specific hostname here: ")
	hostname, _ = reader.ReadString('\n')
	return strings.TrimSuffix(hostname, "\n")
}

func (mach OpenNebulaMachine) GetRepresentative(env Env) string {
	index := GetRolesIndex(mach.Roles)
	if hostname, ok := env.hostnameMapping[index]; ok {
		return hostname
	}
	return ""
}

func (mach OpenNebulaMachine) CommentEtcHosts() error {
	dest := fmt.Sprintf("root@%s", mach.IP)
	cmd := exec.Command("ssh", "-oBatchMode=yes", dest, "echo \"# Opennebula Hosts\" >> /etc/hosts")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (mach OpenNebulaMachine) GetEtcHosts(env Env, wg *sync.WaitGroup) error {
	defer wg.Done()
	dest := fmt.Sprintf("root@%s", mach.IP)
	cmd := exec.Command("ssh", "-oBatchMode=yes", dest, "cat /etc/hosts")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return err
	}
	hosts_file := string(output)

	src := fmt.Sprintf("%s/%s.hosts", env.config.ConfigPath, mach.Name)
	f, err := os.Create(src)
	defer f.Close()
	if err != nil {
		fmt.Println("Error created new host file")
		panic(err)
	}

	hosts_lines := strings.Split(hosts_file, "\n")
	for _, line := range hosts_lines {
		f.WriteString(line + "\n")
		if line == "# Opennebula Hosts" {
			break
		}
	}

	for _, machine := range env.machs {
		if machine.Status && machine.Name != mach.Name {
			line := fmt.Sprintf("%s %s\n", machine.IP, machine.Hostname)
			f.WriteString(line)
		}
	}
	scp_dest := fmt.Sprintf("root@%s:/etc/hosts", mach.IP)

	cmd2 := exec.Command("scp", src, scp_dest)
	err2 := cmd2.Run()
	if err2 != nil {
		fmt.Println("Error in scp of host file")
		panic(err2)
	}
	os.Remove(src)
	return nil
}

func UpdateHostFiles(env Env) error {
	var wg sync.WaitGroup
	for _, mach := range env.machs {
		if mach.Status {
			wg.Add(1)
			go mach.GetEtcHosts(env, &wg)
		}
	}
	wg.Wait()
	return nil
}

type String struct {
	Values []string
}
type Roles struct {
	Roles []string `yaml:"roles"`
}

func ReadHostnameData(data []byte) map[string]Roles {
	roles := make(map[string]Roles)
	err := yaml.Unmarshal(data, &roles)
	if err != nil {
		fmt.Printf("Error unmarshalling machines.yml yaml in %s\n", string(data))
		panic(err)
	}
	return roles

}
func GetRolesIndex(roles []string) string {
	sort.Strings(roles)
	index := strings.Join(roles, ",")
	return index
}

func GetHostnameFile(config *Config) {
	path := config.ConfigPath + "/hostnameMapping.yml"
	if _, err := os.Stat(path); err == nil {
		return
	}
	fileUrl := "http://global.mt.lldns.net/llnwsae/llnwsae-tester/tests/hostnameMapping.yml"
	DownloadFile(path, fileUrl)

}
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
func ReadHostnameFile(config *Config) []byte {
	GetHostnameFile(config)
	path := config.ConfigPath + "/hostnameMapping.yml"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Error unmarshalling hostnameMapping.yml yaml in %s\n", string(data))
		panic(err)
	}
	return data

}
func GetHostnameMapping(config *Config) map[string]string {
	data := ReadHostnameFile(config)
	return GetHostnameMap(data)
}

func GetHostnameMap(data []byte) map[string]string {
	RolesToHostnames := make(map[string]string)
	roles := ReadHostnameData(data)
	for name, role := range roles {
		RolesToHostnames[GetRolesIndex(role.Roles)] = name
	}

	return RolesToHostnames
}
