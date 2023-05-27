package main

import (
	"fmt"
	"os"

	"github.com/robertranjan/kpcli/cmd"
	"github.com/robertranjan/kpcli/versioninfo"
	"github.com/urfave/cli/v2"
)

var Version string

func main() {
	app := cli.NewApp()
	app.Name = "kpcli"
	app.Usage = "kpcli ls --help"
	// NOTE: setting version using below commands
	// 		git rev-parse --short HEAD
	// 		git rev-list HEAD --count
	// app.Version = "2023Feb19.f774b64.10"
	app.Version = versioninfo.Version
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "log-level",
			Usage:   "run in debug mode",
			EnvVars: []string{"KDBX_LOG_LEVEL"},
		},
		&cli.StringFlag{
			Name:    "log-dir",
			Usage:   "location on disk to write logs (optional)",
			EnvVars: []string{"KDBX_LOG_DIR"},
		},
		&cli.StringFlag{
			Name:    "keyfile",
			Usage:   "fullpath of keyfile",
			Aliases: []string{"kf", "k"},
			EnvVars: []string{"KDBX_KEYFILE"},
		},
		&cli.StringFlag{
			Name:    "database",
			Usage:   "kdbx files fullpath",
			Aliases: []string{"db", "dbfile"},
			EnvVars: []string{"KDBX_DATABASE"},
		},
		&cli.StringFlag{
			Name:    "pass",
			Usage:   "kdbx pass",
			Aliases: []string{"p"},
			EnvVars: []string{"KDBX_PASSWORD"},
		},
	}
	app.Commands = []*cli.Command{
		cmd.CmdLs,
		cmd.CmdCreatedb,
		cmd.CmdDiff,
		cmd.CmdAdd,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("   sorry, failed to serve your request: \n\t%v\n", err)
		os.Exit(1)
	}
}
