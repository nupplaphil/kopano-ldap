package main

import (
	"fmt"
	"github.com/nupplaphil/kopano-ldap/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Execute(Version())
}

func Version() string {
	return fmt.Sprintf("version %v\ncommit %v, built at %v", version, commit, date)
}
