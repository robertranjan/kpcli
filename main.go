package main

import (
	"os"

	"github.com/robertranjan/kpcli/cmd"
	"github.com/robertranjan/kpcli/version"
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
	app.Version = version.Version
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
		&cli.BoolFlag{
			Name:    "nokey",
			Usage:   "do not use keyfile - go less secure",
			Aliases: []string{"nk", "n"},
			EnvVars: []string{"KDBX_NOKEY"},
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
		&cli.StringFlag{
			Name:    "config",
			Usage:   "read configs from file",
			Aliases: []string{"c"},
			EnvVars: []string{"KDBX_CONFIG"},
		},
		&cli.StringFlag{
			Name:    "sample-config",
			Usage:   "generate a sample config file: kpcli.toml",
			Aliases: []string{"sample"},
			EnvVars: []string{"KDBX_SAMPLECONFIG"},
		},
	}
	app.Commands = []*cli.Command{
		cmd.CmdLs,
		cmd.CmdCreatedb,
		cmd.CmdDiff,
		cmd.CmdAdd,
		cmd.CmdGenerateSampleConfig,
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
