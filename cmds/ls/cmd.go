package ls

import (
	"log"

	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:  "ls",
	Usage: "lists entries",
	Description: `
	List all entries from a .kdbx database

	Usage

	./kpcli --keyfile <keyfile> \
			--database <database-filename> \
			[--cachefile <cachefile name>] \
		ls  [--reverse] [--days 10] [--sort-by-col 1|2|3|4]
			; --reverse -> reverse order
			; --sort-by-col 1 -> title | 2 -> history count | 3 -> creation time | 4 -> mod time
			; --days 10 --> shows entries created or modified in the last 10 days

		┌─────┬────────────────────────────┬───────────┬─────────────────────┬─────────────────────┐
		|     | 	     (COL: 1)          |  (COL: 2) |        (COL: 3)     |       (COL: 4)      |
		├─────┬────────────────────────────┬───────────┬─────────────────────┬─────────────────────┤
		│     │ TITLE                      │ HISTORIES │ CREATED             │ MODIFIED            │
		├─────┼────────────────────────────┼───────────┼─────────────────────┼─────────────────────┤
		│   1 │ Root/TestEntry             │ 1         │ 2018-02-20 22:58:02 │ 2018-02-21 11:42:45 │
		│   2 │ Recycle Bin/testentry2     │ 7         │ 2018-03-25 18:52:47 │ 2018-03-25 22:24:48 │
		└─────┴────────────────────────────┴───────────┴─────────────────────┴─────────────────────┘

	Example:
		kpcli ls --sortby-col 4 -d 2
			; shows entries modified in last 2 days ORDER by col 4(modified time)

		kpcli \
		--keyfile ./tmp/master-db.key \
		--pass 'super_secret' \
		--db ./tmp/master-db.kdbx \
		ls


	`,

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
		&cli.BoolFlag{
			Name:    "diff",
			Usage:   "diff against cache",
			Aliases: []string{"di"},
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

func runLs(app *cli.Context) error {

	opts := Options{
		Reverse:      app.Bool("reverse"),
		Days:         app.Int("days"),
		Pass:         app.String("pass"),
		Database:     app.String("database"),
		Key:          app.String("keyfile"),
		Sort:         app.String("sort"),
		Fields:       app.String("fields"),
		SortbyCol:    app.Int("sortby-col"),
		CacheFile:    app.String("cachefile"),
		Diff:         app.Bool("diff"),
		Quite:        app.Bool("quite"),
		OutputFormat: app.String("output-format"),
	}

	if opts.Pass == "" {
		log.Fatalf("This command is using:\n"+
			"   database: %v\n"+
			"   keyfile: %v\n"+
			"   \033[33mCould not find password for database.\033[0m\n"+
			"   you may need to \033[32m'source ~/dotfiles/tools/kpcli/.envrc?'\033[0m\n",
			opts.Database, opts.Key)
	}

	ls, err := NewDB(opts)
	if err != nil {
		return err
	}

	return ls.Run()
}
