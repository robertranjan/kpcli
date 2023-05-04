package createdb

import (
	"github.com/urfave/cli/v2"
)

type Options struct {
	Pass     string
	Key      string
	Database string
}

type db struct {
	Options *Options
}

var Cmd = &cli.Command{
	Name:  "createdb",
	Usage: "Create a new kdbx databse",
	Description: `
	To list entries

	kpcli \
		--keyfile <keyfile> \
		--name <xyx.kdbx> \
		--pass {password to encrypt/open kdbx} \
		createdb

	Example:

		kpcli \
			--keyfile ./tmp/master-db.key \
			--pass 'super_secret' \
			--db ./tmp/master-db.kdbx \
			createdb
	`,

	Action: runCreate,
	Flags:  []cli.Flag{
		// not using this sort stringFlag option, yet
	},
}

func runCreate(app *cli.Context) error {

	opts := Options{
		Database: app.String("database"),
		Pass:     app.String("pass"),
		Key:      app.String("keyfile"),
	}
	db := &db{
		Options: &opts,
	}

	return db.Run()
}
