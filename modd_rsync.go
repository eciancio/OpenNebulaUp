package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/shirou/gopsutil/process"
)

type ModdPid struct {
	Pid int32 `json:"pid"`
}

// Start the modd proccess with modd conf file
func StartModd(env Env) {
	filename := "opennebula_modd-" + env.environment + ".conf"
	cmd := exec.Command("modd", "-f", filename)
	cmd.Dir = env.config.ConfigPath
	errs := cmd.Start()
	if errs != nil {
		fmt.Println("Unable to run: ", errs)
		return
	}
	WriteModdPid(cmd.Process.Pid, env)
}

func StopModd(env Env) {
	pid := GetModdPid(env)
	p, err := process.NewProcess(pid)
	if err != nil {
		fmt.Println(err)
		return
	}

	name, err := p.Name()
	if err != nil {
		fmt.Println(err)
		return
	}
	if name != "modd" {
		return
	}
	process, err := os.FindProcess(int(pid))
	if err != nil {
		fmt.Printf("ERR: Process pid: %d could not be found\n", pid)
		return
	}
	err = process.Kill()
	if err != nil {
		fmt.Println(err)
	}
	moddPIDPath := GetModdPIDPath(env)
	os.Remove(moddPIDPath) // Remove when done
}

func GetModdPIDPath(env Env) string {
	moddPidName := env.config.ConfigPath + "/modd_pid-" + env.environment + ".json"
	return moddPidName
}
func WriteModdPid(pid int, env Env) {
	modd_pid := ModdPid{int32(pid)}

	file, _ := json.MarshalIndent(modd_pid, "", " ")
	moddPIDPath := GetModdPIDPath(env)

	_ = ioutil.WriteFile(moddPIDPath, file, 0644)
}

func GetModdPid(env Env) int32 {
	var configuration ModdPid
	moddPIDPath := GetModdPIDPath(env)
	configFile, err := os.Open(moddPIDPath)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&configuration)
	return configuration.Pid
}
