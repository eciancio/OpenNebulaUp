package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"syscall"

	"github.com/urfave/cli"
)

type Env struct {
	connection      gocaConnection
	machs           map[string]OpenNebulaMachine
	arguements      []string
	data            ioData
	hostnameMapping map[string]string
	config          *Config
	environment     string
}

func main() {
	data := fileIO{}
	arguements := os.Args
	ConfigPath := GetConfigLocation(arguements)
	config, err := GetConfig(ConfigPath)
	config.ConfigPath = GetUserHomeDir() + "/.config/OpenNebulaUp/"
	if err != nil {
		return
	}

	ProjectsPath := GetAltProjectsPath(arguements)
	if ProjectsPath != "" {
		config.ProjectsPath = ProjectsPath
	}
	// Go To projects path
	err = os.Chdir(config.ProjectsPath + "salt-states")
	if err != nil {
		fmt.Printf("WARNING: salt-states was not found in %s. Please check if %s is correct\n", config.ProjectsPath, config.ProjectsPath)
	}

	c := gocaConnection{GetOpenNebulaConfig(config)}
	machs := GetAllMachines(&data)
	// machs = SetMachinesStatus(&c, machs)
	hostnameMapping := GetHostnameMapping(config)
	env := Env{c, machs, arguements, &data, hostnameMapping, config, "base"}

	app := GetAppConfig(env)
	app.Run(os.Args)
}

func GetAppConfig(env Env) *cli.App {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "",
			Usage: "Path to alternative config default is $HOME/.config/OpenNebulaUp/opennebulaup.config",
		},
		cli.StringFlag{
			Name:        "env",
			Value:       "base",
			Usage:       "The OpenNebulaEnvironment",
			Destination: &env.environment,
		},
		cli.StringFlag{
			Name:        "ProjectsPath",
			Value:       env.config.ProjectsPath,
			Usage:       "An Alternative Projects Path not defined in config",
			Destination: &env.config.ProjectsPath,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "up",
			Usage: "starts and provisions the opennebula machine",
			Action: func(c *cli.Context) error {
				mach := c.Args().First()
				if mach == "" {
					fmt.Println("Please Specify a machine")
					return nil
				}

				env.machs = SetMachinesStatus(&env.connection, env.machs, env)
				OpenNebulaUp(env, mach, false)
				return nil
			},
		},
		{
			Name:  "hold",
			Usage: "starts the opennebula machine on hold ( recommended to use up )",
			Action: func(c *cli.Context) error {
				mach := c.Args().First()
				if mach == "" {
					fmt.Println("Please Specify a machine")
					return nil
				}

				env.machs = SetMachinesStatus(&env.connection, env.machs, env)
				OpenNebulaHold(env, mach)
				return nil
			},
		},
		{
			Name:  "list",
			Usage: "list all the available machines without status",
			Action: func(c *cli.Context) error {
				OpenNebulaList(env.machs)
				return nil
			},
		},

		{
			Name:  "destroy",
			Usage: "stops and deletes all traces of the opennebula  machine",
			Action: func(c *cli.Context) error {
				mach := c.Args().First()
				env.machs = SetMachinesStatus(&env.connection, env.machs, env)
				if mach == "" {
					OpenNebulaDestroyAll(env, c.Bool("f"))
					return nil
				}

				OpenNebulaDestroy(env, mach)
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "f, force",
				},
			},
		},
		{
			Name:  "resume",
			Usage: "Resume a VM that is in poweroff or suspend mode",
			Action: func(c *cli.Context) error {
				mach := c.Args().First()
				if mach == "" {
					fmt.Println("Please Specify a machine")
					return nil
				}

				env.machs = SetMachinesStatus(&env.connection, env.machs, env)
				OpenNebulaResume(env, mach)
				return nil
			},
		},
		{
			Name:  "reboot",
			Usage: "Reboot a VM that is currently in OpenNebula",
			Action: func(c *cli.Context) error {
				mach := c.Args().First()
				if mach == "" {
					fmt.Println("Please Specify a machine")
					return nil
				}

				env.machs = SetMachinesStatus(&env.connection, env.machs, env)
				OpenNebulaReboot(env, mach)
				return nil
			},
		},
		{
			Name:  "poweroff",
			Usage: "Turn off VM but keep allocated resources",
			Action: func(c *cli.Context) error {
				mach := c.Args().First()
				if mach == "" {
					fmt.Println("Please Specify a machine")
					return nil
				}

				env.machs = SetMachinesStatus(&env.connection, env.machs, env)
				OpenNebulaPowerOff(env, mach)
				return nil
			},
		},

		{
			Name:  "created",
			Usage: "outputs only cretaed machines",
			Action: func(c *cli.Context) error {
				env.machs = SetMachinesStatus(&env.connection, env.machs, env)
				PrintCreatedStatus(env.machs)
				return nil
			},
		},

		{
			Name:  "status",
			Usage: "outputs status of the opennebula machines",
			Action: func(c *cli.Context) error {
				mach := c.Args().First()
				env.machs = SetMachinesStatus(&env.connection, env.machs, env)
				if mach == "" {
					PrintAllStatus(env.machs)
					return nil
				}

				PrintSingleStatus(env.machs, mach)
				return nil
			},
		},
		{
			Name:  "rsync",
			Usage: "Restarts Rsync watcher for salt-master ",
			Action: func(c *cli.Context) error {
				env.machs = SetMachinesStatus(&env.connection, env.machs, env)
				RestartRsync(env)
				return nil
			},
		},
		{
			Name:  "print",
			Usage: "prints all relevant data for particular machine",
			Action: func(c *cli.Context) error {
				mach := c.Args().First()
				if mach == "" {
					fmt.Println("Please Specify A Machine")
					return nil
				}

				PrintSingleMachine(env, mach)
				return nil
			},
		},
		{
			Name:  "ssh",
			Usage: "connects to machine via SSH",
			Action: func(c *cli.Context) error {
				mach := c.Args().First()
				if mach == "" {
					fmt.Println("Please Specify a machine")
					return nil
				}
				env.machs = SetMachinesStatus(&env.connection, env.machs, env)
				OpenNebulaSsh(env, mach)

				return nil
			},
		},
		{
			Name:   "completion",
			Usage:  "output bash completion script. source <(OpenNebulaUp completion)",
			Action: completion,
		},
	}
	app.Version = "0.2.1"
	app.Name = "OpenNebulaUp"
	app.Usage = "CLI for managing OpenNebula VMs"

	return app
}

func OpenNebulaUp(env Env, mach string, hold bool) error {
	IP := GetMasterIP(env)
	if IP == "0" && mach != "salt-master" {
		fmt.Println("Warning salt-master is not started")
		fmt.Println("Are you sure you want to continue (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if text != "y\n" {
			os.Exit(1)
		}
	}
	if _, ok := env.machs[mach]; !ok {
		fmt.Printf("%s is not defined in any machines.yml file\n", mach)
		os.Exit(1)
	}
	machine := env.machs[mach].StartOpenNebulaMachine(env, IP, hold)
	if !hold {
		machine.WaitForBoot()
		machine.CommentEtcHosts()
		env.machs[mach] = machine
		UpdateHostFiles(env)
		if mach == "salt-master" {
			CreateModdConf(machine.IP, env)
			StartModd(env)
			machine.WaitForSync(env)
		}

	}

	return nil
}

func OpenNebulaHold(env Env, mach string) error {
	OpenNebulaUp(env, mach, true)
	return nil
}

func OpenNebulaDestroyAll(env Env, force bool) error {
	if !force {
		fmt.Println("Are you sure you want to destroy all (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if text != "y\n" {
			return nil
		}
	}
	for _, mach := range env.machs {
		if mach.Status {
			mach.DestroyOpenNebulaMachine(env)
		}
	}
	return nil
}

func OpenNebulaDestroy(env Env, mach string) error {
	env.machs[mach].DestroyOpenNebulaMachine(env)
	return nil
}

func OpenNebulaReboot(env Env, mach string) error {
	env.machs[mach].RebootOpenNebulaMachine(env)
	return nil
}
func OpenNebulaResume(env Env, mach string) error {
	env.machs[mach].ResumeOpenNebulaMachine(env)
	return nil
}

func OpenNebulaPowerOff(env Env, mach string) error {
	env.machs[mach].PowerOffOpenNebulaMachine(env)
	return nil
}

func OpenNebulaSsh(env Env, mach string) error {
	path, err := exec.LookPath("ssh")
	if err != nil {
		return err
	}
	IP := env.machs[mach].IP
	ip_string := "root@" + IP

	syscall.Exec(path, []string{"ssh", "-oBatchMode=yes", ip_string}, []string{"TERM=screen-256color"})
	return nil
}

func RestartRsync(env Env) {
	if master, ok := env.machs["salt-master"]; ok && master.Status {
		StopModd(env)
		CreateModdConf(master.IP, env)
		StartModd(env)
		fmt.Println("Rsync with master restarted")
		return

	}
	fmt.Println("Can't restart Rsync salt-master not running")
}

func GetMasterIP(env Env) string {
	if env.machs["salt-master"].Status {
		return env.machs["salt-master"].IP
	}
	return "0"
}
func PrintAllStatus(machs map[string]OpenNebulaMachine) {
	keys := make([]string, 0)
	for k, _ := range machs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		status := "not created (OpenNebula)"
		mach := machs[key]
		if mach.Status {
			status = "Running (%s)"
			status = fmt.Sprintf(status, mach.IP)
		}

		fmt.Printf("%-30s %-s\n", mach.Name, status)
	}
}

func PrintSingleStatus(machs map[string]OpenNebulaMachine, name string) {
	status := "not created (OpenNebula)"
	mach := machs[name]
	if mach.Status {
		status = "Running (%s)"
		status = fmt.Sprintf(status, mach.IP)
	}

	fmt.Printf("%-30s %-s\n", mach.Name, status)
}
func PrintCreatedStatus(machs map[string]OpenNebulaMachine) {
	var status string
	for _, mach := range machs {
		if mach.Status {
			status = "Running (%s)"
			status = fmt.Sprintf(status, mach.IP)
			fmt.Printf("%-30s %-s\n", mach.Name, status)
		}
	}

}

func OpenNebulaList(machs map[string]OpenNebulaMachine) {
	for _, mach := range machs {
		fmt.Printf("%s\n", mach.Name)
	}
}

func PrintSingleMachine(env Env, name string) {
	mach := env.machs[name]
	fmt.Printf("Machine Name: %s\n", name)
	hostname := mach.ReadHostname(env)
	fmt.Printf("Hostname: %s\n", hostname)
	fmt.Printf("Operating System: %s\n", mach.OperatingSystem)
	fmt.Printf("Mem: %d\n", mach.Mem)
	fmt.Printf("VCPU: %d\n", mach.VCPU)
	fmt.Printf("Team Environment: %s\n", mach.TeamEnv)
	fmt.Println("Roles :")
	for _, role := range mach.Roles {
		fmt.Printf("   - %s\n", role)
	}

}

type Machine struct {
	Inherit   string `yaml:"inherit"`
	Enabled   bool   `yaml:"enabled"`
	Autostart bool   `yaml:"autostart"`
	Mem       int    `yaml:"mem"`
	Hostname  string `yaml:"hostname"`
	VCPU      int    `yaml:"vcpu"`
	TeamEnv   string `yaml:"team_env"`
}
type Tags struct {
	Tags []map[string]string
}

type OpenNebulaMachine struct {
	Name            string
	OperatingSystem string
	Roles           []string
	IP              string
	ID              uint
	Status          bool
	Hostname        string
	Mem             int
	VCPU            int
	TeamEnv         string
	TagData         Tags
}

func GetMachineObjects() map[string]Machine {
	return make(map[string]Machine)
}
