package main

import (
	"fmt"
	"testing"
)

var expected_output = "testing/path/** {\n    prep: rsync -a -e \"ssh\" --rsync-path=\"sudo rsync\" testing/path/salt-states/ root\\@1:/srv/teams/salt-states/\n    prep: rsync -a -e \"ssh\" --rsync-path=\"sudo rsync\" testing/path/salt-sae/ root\\@1:/srv/teams/salt-sae/\n    prep: rsync -a -e \"ssh\" --rsync-path=\"sudo rsync\" testing/path/salt-pillar/ root\\@1:/srv/pillar-git/\n    prep: rsync -a -e \"ssh\" --rsync-path=\"sudo rsync\" opennebula_modd-base.conf root\\@1:/srv/\n}"

var expected_output2 = "path/to/home/** {\n    prep: rsync -a -e \"ssh\" --rsync-path=\"sudo rsync\" path/to/home/salt-states/ root\\@2:/srv/teams/salt-states/\n    prep: rsync -a -e \"ssh\" --rsync-path=\"sudo rsync\" path/to/home/salt-sae/ root\\@2:/srv/teams/salt-sae/\n    prep: rsync -a -e \"ssh\" --rsync-path=\"sudo rsync\" path/to/home/salt-pillar/ root\\@2:/srv/pillar-git/\n    prep: rsync -a -e \"ssh\" --rsync-path=\"sudo rsync\" opennebula_modd-base.conf root\\@2:/srv/\n}"

func TestCreateModdConf(t *testing.T) {
	env := Env{}
	env.environment = "base"
	config := Config{}
	config.ProjectsPath = "testing/path"
	data := MockIO{}
	env.data = &data
	env.config = &config
	CreateModdConf("1", env)
	if expected_output != data.data {
		fmt.Println(data.data)
		t.Errorf("Output of moddConf is not expected")
	}

}
func TestCreateModdConf2(t *testing.T) {
	env := Env{}
	env.environment = "base"
	config := Config{}
	config.ProjectsPath = "path/to/home"
	data := MockIO{}
	env.data = &data
	env.config = &config
	CreateModdConf("2", env)
	if expected_output2 != data.data {
		t.Errorf("Output of moddConf is not expected")
	}

}
