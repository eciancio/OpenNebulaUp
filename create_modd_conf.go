package main

import (
	"fmt"
	"os"
	"strings"
)

func CreateModdConf(IP string, env Env) {

	if env.config.ProjectsPath[len(env.config.ProjectsPath)-1:] != "/" {
		env.config.ProjectsPath += "/" // Ensure that Path Ends with /
	}
	err := os.Chdir(env.config.ProjectsPath + "salt-states") // If path has no salt-states give warning but continue
	if err != nil {
		fmt.Printf("WARNING: salt-states was not found in %s. Please check if %s is correct\n", env.config.ProjectsPath, env.config.ProjectsPath)
	}

	modd_path := env.config.ConfigPath + "opennebula_modd-" + env.environment + ".conf" // Filename is environemnt specific
	data := env.data
	data.SetDataLocation(modd_path)
	data.Create()
	parent := env.config.ProjectsPath
	if strings.HasSuffix(parent, "/") {
		parent = parent[:len(parent)-1]
	}

	watch_line := parent + "/** {\n"
	data.WriteString(watch_line)
	pillar_string := "    prep: rsync -a -e \"ssh\" --rsync-path=\"sudo rsync\" %s/%s/ root\\@%s:/srv/pillar-git/\n"
	base_string := "    prep: rsync -a -e \"ssh\" --rsync-path=\"sudo rsync\" %s/%s/ root\\@%s:/srv/teams/%s/\n"
	repos := FindAllSaltDirectories(env.data)
	data.WriteString(fmt.Sprintf(base_string, parent, "salt-states", IP, "salt-states"))
	for _, repo := range repos {
		if repo == "salt-pillar" {
			line := fmt.Sprintf(pillar_string, parent, repo, IP)
			data.WriteString(line)
		} else {
			line := fmt.Sprintf(base_string, parent, repo, IP, repo)
			data.WriteString(line)
		}
	}

	confString := "    prep: rsync -a -e \"ssh\" --rsync-path=\"sudo rsync\" %s root\\@%s:/srv/\n"
	line := fmt.Sprintf(confString, modd_path, IP)
	data.WriteString(line)
	data.WriteString("}\n")

	synced_string := "    prep: rsync -a -e \"ssh\" --rsync-path=\"sudo mkdir -p %s && sudo rsync\" %s root\\@%s:%s/\n"
	for _, mach := range env.machs {
		if len(mach.SyncedFolders) > 0 && mach.Status {

			for _, folder := range mach.SyncedFolders {
				parent = folder.Source
				if strings.HasSuffix(parent, "/") {
					parent = parent[:len(parent)-1]
				}
				watch_line := parent + "/** {\n"
				data.WriteString("\n")
				data.WriteString(watch_line)
				line := fmt.Sprintf(synced_string, folder.Destination, folder.Source, mach.IP, folder.Destination)
				data.WriteString(line)
				data.WriteString("}\n")
			}
		}
	}

}
