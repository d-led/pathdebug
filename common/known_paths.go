package common

import (
	"github.com/mitchellh/go-homedir"
)

var knownPaths = []string{
	// global
	"/etc/profile",
	"/etc/zprofile",
	"/etc/zlogin",
	"/etc/zshenv",
	// generic
	"~/.profile",
	// zsh: https://zsh.sourceforge.io/Doc/Release/Files.html#Startup_002fShutdown-Files
	"~/.zshenv",
	"~/.zprofile",
	"~/.zshrc",
	"~/.zlogin",
	// bash: https://www.gnu.org/software/bash/manual/html_node/Bash-Startup-Files.html
	"~/.bashrc",
	"~/.bash_profile",
	"~/.bash_login",
}

// ForEachKnownPath iterates over expanded known paths if they're legitimate
func ForEachKnownPath(fn func(string, string)) {
	for _, knownPath := range knownPaths {
		expanded, err := homedir.Expand(knownPath)
		if err != nil {
			continue
		}
		fn(knownPath, expanded)
	}
}
