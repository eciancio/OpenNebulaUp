package main

import (
	"fmt"
	"os/exec"
	"time"
)

func (mach OpenNebulaMachine) WaitForBoot() error {
	fmt.Printf("%s is booting\n", mach.Name)
	time.Sleep(10 * time.Second)
	timeoutNumber := 10
	booted := false
	for i := 0; i < timeoutNumber; i++ {
		booted = mach.HasBooted()
		if booted {
			break
		}
		time.Sleep(10 * time.Second)
	}
	if !booted {
		fmt.Printf("%s timed out on booting", mach.Name)
		return nil
	}

	fmt.Printf("%s finished  booting\n", mach.Name)
	return nil
}

func (mach OpenNebulaMachine) HasBooted() bool {
	dest := fmt.Sprintf("root@%s", mach.IP)
	cmd := exec.Command("ssh", "-oBatchMode=yes", dest, "ls")
	err := cmd.Run()
	if err != nil {
		return false
	}

	return true
}

func ClearKeys(env Env) error {
	if !env.machs["salt-master"].Status {
		fmt.Println("Can't clear keys salt-master is not running")
		return nil
	}

	masterIP := env.machs["salt-master"].IP
	dest := fmt.Sprintf("root@%s", masterIP)
	cmd := exec.Command("ssh", "-oBatchMode=yes", dest, "salt-key", "-D", "-y")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (mach OpenNebulaMachine) WaitForSync(env Env) error {
	fmt.Printf("%s is sycning salt\n", mach.Name)
	time.Sleep(15 * time.Second)
	timeoutNumber := 20
	synced := false
	for i := 0; i < timeoutNumber; i++ {
		synced = CheckSyncFinished(mach, env)
		if synced {
			break
		}
		time.Sleep(10 * time.Second)
	}
	if !synced {
		fmt.Printf("%s timed out on sycning salt repos", mach.Name)
		return nil
	}

	fmt.Printf("%s finished syncing\n", mach.Name)
	return nil

}
func CheckSyncFinished(mach OpenNebulaMachine, env Env) bool {
	dest := fmt.Sprintf("root@%s", mach.IP)
	name := "ls -l /srv/opennebula_modd-" + env.environment + ".conf"
	cmd := exec.Command("ssh", "-oBatchMode=yes", dest, name)
	err := cmd.Run()
	if err != nil {
		return false
	}

	return true
}
