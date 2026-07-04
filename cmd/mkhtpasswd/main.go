package main

import (
	"fmt"
	"os"

	"mkhtpasswd/internal/cli"
	"mkhtpasswd/internal/ui"
)

func main() {
	cfg, ok := cli.Parse(os.Args[1:])
	if !ok {
		os.Exit(2)
	}

	if cfg.ShowVersion || (!cfg.ForceGUI && hasCLIArgs(cfg)) {
		res, err := cli.Run(cfg)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		cli.ShowResult(cfg, res)
		return
	}

	if err := ui.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func hasCLIArgs(cfg cli.Config) bool {
	return cfg.User != "" || cfg.Random || cfg.Password != "" || cfg.ForceGUI
}