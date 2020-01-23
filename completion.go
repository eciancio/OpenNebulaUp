package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func completion(c *cli.Context) error {
	out := `#!/bin/bash

_active() {
    OpenNebulaUp created | sed -E 's/\s+.*$//' | tr '\n' ' ' | sed -E 's/\s+$//'
}

_all() {
    OpenNebulaUp list| tr '\n' ' ' | sed -E 's/\s+$//'
}

_one_complete() {
    local cur prev active cmds
    cmds="completion up hold destroy list resume reboot poweroff created status clear-keys csmock rsync print ssh help --version -v --help -h --config"
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    case "${prev}" in
        ssh|hold|poweroff|reboot|resume|destroy)
            active=$(_active)
            COMPREPLY=($(compgen -W "${active}" -- "${cur}"))
            return 0
            ;;
        up|print)
            active=$(_all)
            COMPREPLY=($(compgen -W "${active}" -- "${cur}"))
            return 0
            ;;
        clear-keys|help|status|rsync|created|csmock|list)
            return 0
            ;;
        --config)
            compopt -o default; COMPREPLY=()
            return 0
            ;;
        *)
        ;;
    esac

   COMPREPLY=($(compgen -W "${cmds}" -- ${cur}))
   return 0
}

complete -F _one_complete OpenNebulaUp
`
	fmt.Print(out)
	return nil
}
