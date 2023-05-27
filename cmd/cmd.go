package cmd

import (
	"github.com/urfave/cli/v2"
)

var CmdAdd = &cli.Command{
	Name:    "add",
	Usage:   "add an entry",
	Aliases: []string{"l"},
	Description: `Add an entry to a .kdbx database

syntax:
	./kpcli --keyfile <keyfile> \
			--database <database-filename> \
			--pass <pass to open database> \
		add --title <title> \
			--user <username> \
			--pass <password>

Example:
		kpcli add \
			--title new-entry-1 \
			--user example-user1 \
			--pass secret_13

	`,

	Action: runAddEntry,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "entry-title",
			Usage:   "title of new entry",
			Aliases: []string{"t"},
		},
		&cli.StringFlag{
			Name:    "entry-user",
			Usage:   "user of new entry",
			Aliases: []string{"u"},
		},
		&cli.StringFlag{
			Name:    "entry-pass",
			Usage:   "pass of new entry",
			Aliases: []string{"p"},
		},
	},
}

var CmdCreatedb = &cli.Command{
	Name:    "create",
	Usage:   "Create a new kdbx databse",
	Aliases: []string{"c"},
	Description: `create command create a new kdbx database with few sample entries

Example:
	kpcli \
		--keyfile ./tmp/master-db.key \
		--pass 'super_secret' \
		--db ./tmp/master-db.kdbx \
		create`,

	Action: runCreate,
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "entries",
			Usage:   "number of sample entries",
			Aliases: []string{"e"},
		},
	},
}

var CmdDiff = &cli.Command{
	Name:    "diff",
	Usage:   "diff entries between 2 kdbx databases",
	Aliases: []string{"d"},
	Description: `Show difference between 2 kdbx databases

syntax:
	kpcli \
		--keyfile <keyfile> \
		--database <database-filename> \
		--pass "${KDBX_PASSWORD}" \
		diff \
			--database2 <database-filename-2>

example:
	kpcli \
		--keyfile ${KDBX_KEYFILE} \
		--database ${KDBX_PASSWORD} \
		--pass "${KDBX_PASSWORD}" \
		diff \
			--database2 ${DATABASE_BACKUP}
	`,

	Action: runDiff,
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

var CmdLs = &cli.Command{
	Name:    "ls",
	Usage:   "lists entries",
	Aliases: []string{"l"},
	Description: `List all entries from a .kdbx database

Example:
	kpcli \
		--keyfile ./tmp/master-db.key \
		--pass 'super_secret' \
		--db ./tmp/master-db.kdbx \
		ls`,

	Action: runLs,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "fields",
			Usage:   "fields list to be displayed",
			Aliases: []string{"f"},
		},
		&cli.BoolFlag{
			Name:    "reverse",
			Usage:   "in reverse order",
			Aliases: []string{"r"},
		},
		&cli.BoolFlag{
			Name:    "quite",
			Usage:   "less verbose",
			Aliases: []string{"q"},
		},
		&cli.StringFlag{
			Name:    "cachefile",
			Usage:   "cache result",
			Aliases: []string{"ca"},
		},
		&cli.StringFlag{
			Name:    "days",
			Usage:   "number of days ; days <= 0 means all",
			Aliases: []string{"d"},
		},
		&cli.StringFlag{
			Name:    "sortby-col",
			Usage:   "sort by column number starting 1",
			Aliases: []string{"sb"},
		},
		&cli.StringFlag{
			Name:    "output-format",
			Usage:   "Output format; available: table, csv, markdown, html",
			Aliases: []string{"of"},
		},
	},
}
