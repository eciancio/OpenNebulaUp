package main

import (
	"fmt"

	"github.com/urfave/cli"
)

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
