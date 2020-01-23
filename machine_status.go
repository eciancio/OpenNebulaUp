package main

import (
	"github.com/OpenNebula/one/src/oca/go/src/goca"
)

type Connector interface {
	GetAllVms() *goca.VMPool
	GetVMIP(uint) string
	StartVM(string, string, uint, bool) (string, uint)
	DestroyVM(uint) error
}

type gocaConnection struct {
	config goca.OneConfig
}

func (con *gocaConnection) GetAllVms() *goca.VMPool {
	goca.SetClient(con.config)
	pool, err := goca.NewVMPool(-3)
	if err != nil {
		panic(err)
	}
	return pool
}

func (con *gocaConnection) GetVMIP(ID uint) string {
	goca.SetClient(con.config)
	minion := goca.NewVM(ID)
	err := minion.Info()
	if err != nil {
		panic(err)
	}

	IP := minion.Template.NICs[0].IP
	return IP
}

func GetAllRunningMachines(c Connector) *goca.VMPool {
	return c.GetAllVms()
}

func GetVMIP(c Connector, ID uint) string {
	return c.GetVMIP(ID)
}

func SetMachinesStatus(c Connector, machs map[string]OpenNebulaMachine, env Env) map[string]OpenNebulaMachine {
	vms := c.GetAllVms()
	for i := 0; i < len(vms.VMs); i++ {
		name := vms.VMs[i].Name
		if mach, ok := machs[name]; ok {
			vms.VMs[i].Info()
			if vms.VMs[i].Template.Context != nil {
				if environment, ok := vms.VMs[i].Template.Context.Dynamic["OPEN_NEBULA_UP_ENV"]; ok && environment == env.environment {
					mach.Status = true

					if host, ok := vms.VMs[i].Template.Context.Dynamic["SET_HOSTNAME"]; ok {
						mach.Hostname = host
					} else {
						mach.Hostname = mach.Name + ".devnet.llnw.net"
					}
					if len(vms.VMs[i].Template.NICs) == 0 {
						mach.IP = "0"
					} else {
						mach.IP = vms.VMs[i].Template.NICs[0].IP
					}
					mach.ID = vms.VMs[i].ID
					machs[name] = mach
				}
			}
		}
	}
	return machs
}
