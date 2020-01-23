# OpenNebulaUp
Command Line Tool for maintaining a virtual development environment with OpenNebula Virtual Machines

## Install 

### Manual
1. Install Modd <br>

2. Install OpenNebulaUp <br>

## Getting Started

#### Set Up Environment
  
3. Create opennebulaup.config in ~/.config/OpenNebulaUp/
  ```
  OpenNebulaToken=""
  OpenNebulaUsername=""
  SshKey=""
  OpenNebulaApi=""
  ProjectsPath=""
  ```
  
 Projects path is the path to your projects folder with all your salt-repos
  
4. Get OpenNebula token from 
  * Login and Naviagte to settings > auth > manage login tokens
  * Create a new token. Set Expiration, in seconds to -1 for non expiring token
  * Copy your token into the token section of ~/.config/OpenNebulaUp/opennebulaup.config
  
5. Add ssh public key to public_key section ~/.config/OpenNebulaUp/opennebulaup.config
  
## Usage 

OpenNebulaUp gets machines definitions from machines.yml and csmock.yml in team repos

Example of machines.yml
```
sae-nexus:
  inherit: salt-xenial
  hostname: nexus-01.phx3
  mem: 4098
  vcpu: 2
  team_env: dev
```

Example of csmock.yml
```
sae-nexus:
  name: sae-nexus
  tags:
    - name: role-sae-nexus
      tag_group_name: vagrant
```
### Commands 
OpenNebulaUp print \<machineName\> - see machines data before starting 
```
$ OpenNebulaUp print sae-nexus
Machine Name: sae-nexus
Hostname: nexus-01.phx3
Operating System: salt-xenial
Mem: 4098
VCPU: 2
Team Environment: dev
Roles :
   - sae-nexus
  ```
  
  OpenNebulaUp status - see status of machines 
  
```
$ OpenNebulaUP status
sae-nexus                      running   (10.12.161.18)
sae-singel                     not created (OpenNebula)
sae-slb                        not created (OpenNebula)
sae-splunk-deploy              not created (OpenNebula)
sae-splunk-indexer1            not created (OpenNebula)
...
```
OpenNebulaUp up \<machinName\> - start machine in openenbula
```
$ OpenNebulaUp up salt-master
salt-master is booting
salt-master finished booting
```

OpenNebulaUp destroy \<machinName\> - will terminate machine in opennebula
```
$ OpenNebulaUp destroy salt-master
Destroying salt-master
```

OpenNebulaUp ssh \<machinName\> - will start an ssh session to the machine as root user
```
$ OpenNebulaUp ssh salt-master
Welcome to Ubuntu 14.04.5 LTS (GNU/Linux 4.4.0-31-generic x86_64)

 * Documentation:  https://help.ubuntu.com/
root@salt-master:~#
```

`OpenNebulaUp -h` - help for additional commands

### Environments
OpenNebulaUp supports multiple OpenNebulaUp Environments. This allows you to have multiple instances of salt-masters and minions in OpenNebula. This also allows for different salt projects directoires to be synced to the master.

```
$ OpenNebula up --env feature1 up salt-master
```
Using --env feature1 will start a new salt-master in the OpenNebulaUp feature1 environment

```
$ OpenNebula up --env feature1 up sae-nexus
```
Using --env feature1 on a minion will connect it the feature1 salt-master

```
$ OpenNebula up --env feature1 created

salt-master                    running   (10.12.161.17)
sae-nexus                      running   (10.12.161.18)
```
You can see all the VMs in the feature1 environment

```
$ OpenNebula up --env feature1 --ProjectsPath ~/LLNW/SaltProjects2/ up salt-master
```
Using --ProjectsPath you can specify a different projects directory to sync to your salt-master. The defualt ProjectsPath is is the one listed in the config.  

### Bash Completion
`OpenNebulaUp completion` will output bash completions. On Linux or MacOS, you can put the following line in your `.bashrc`:

```bash
. <(OpenNebulaUp completion)
```


