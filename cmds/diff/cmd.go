package diff

import (
	"log"
	"path"
	"strings"

	"github.com/urfave/cli/v2"
)

const BackupDIR = "/tmp/backupDIR"

var Cmd = &cli.Command{
	Name:  "diff",
	Usage: "diff entries from 2 kdbx files",
	Description: `
Difference between 2 kdbx databases

syntax:
	kpcli \
		--keyfile <keyfile> \
		--database <database-filename> \
		--pass "${PASS}" \
		diff \
			--database2 <database-filename-2>
			; --database2 - latest file from backup path will be used when it's empty

example:
	kpcli \
		--keyfile ${KEYFILE} \
		--database ${DATABASE} \
		--pass "${PASS}" \
		diff \
			--database2 ${DATABASE_BACKUP}
	`,

	Action: cmd,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "database2",
			Usage:   "kdbx files fullpath2",
			Aliases: []string{"db2", "dbfile2"},
			EnvVars: []string{"KDBX_DATABASE2"},
		},
		// not using this sort stringFlag option, yet
		&cli.StringFlag{
			Name:    "output-format",
			Usage:   "Output format; available: table, csv, markdown, html",
			Aliases: []string{"of2"},
		},
		&cli.BoolFlag{
			Name:    "notify",
			Usage:   "notify with email",
			Aliases: []string{"n"},
		},
	},
}

func cmd(app *cli.Context) error {

	opts := Options{
		Pass:           app.String("pass"),
		Database:       app.String("database"),
		Key:            app.String("keyfile"),
		Database2:      app.String("database2"),
		Notify:         app.Bool("notify"),
		OutputFormat:   "csv",
		OutputFilename: "diffLog2Email.html",
	}

	pattern := strings.Split(path.Base(opts.Database), ".")[0]
	if opts.Database2 == "" {
		opts.Database2 = getRecentFile(BackupDIR, pattern)
	}

	if opts.Pass == "" {
		log.Fatalf("Using:\n"+
			"   database: %v\n"+
			"   database2: %v\n"+
			"   keyfile: %v\n"+
			"   \033[33mCould not find password for database.\033[0m\n"+
			"   you may need to \033[32m'source ~/dotfiles/tools/kpcli/.envrc.%v?'\033[0m\n",
			opts.Database, opts.Database2, opts.Key, pattern)
	}
	diff := NewDiff(opts)
	return diff.Run()
}