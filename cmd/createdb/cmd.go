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
	Name:    "createdb",
	Usage:   "Create a new kdbx databse",
	Aliases: []string{"c"},
	Description: `createdb command create a new kdbx database with few sample entries

syntax:
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
	Flags:  []cli.Flag{},
}

func runCreate(app *cli.Context) error {

	opts := Options{
		Database: app.String("database"),
		Pass:     app.String("pass"),
		Key:      app.String("keyfile"),
	}
	_, err := NewDB(opts)
	if err != nil {
		return err
	}

	return nil
}
