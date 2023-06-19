package cmd

import (
	"fmt"
	"math/rand"
	"path"
	"path/filepath"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/urfave/cli/v2"
)

func runAddEntry(app *cli.Context) error {

	d, err := newObject(app)
	if err != nil {
		fmt.Printf("failed to create db : %v\n", err)
		return err
	}

	if d.Options.Pass == "" {
		d.Options.Pass = "super_secret"
	}

	return d.AddEntry()
}

// func (d *db) setupDefaultKeyKdbx() {

// 	if d.Options.Key == "" {
// 		d.Options.Key = "./tmp/master-db.key"
// 	}
// 	if d.Options.Database == "" {
// 		d.Options.Database = "./tmp/master-db.key"
// 	}
// }

func newObject(app *cli.Context) (*db, error) {

	opts := Options{
		CacheFile:      app.String("cachefile"),
		Database:       app.String("database"),
		Database2:      app.String("database2"),
		Days:           app.Int("days"),
		DiffCalling:    app.Bool("diff-calling"),
		EntryPass:      app.String("entry-pass"),
		EntryTitle:     app.String("entry-title"),
		EntryUser:      app.String("entry-user"),
		Fields:         app.String("fields"),
		Key:            app.String("keyfile"),
		Key2:           app.String("keyfile2"),
		LogLevel:       app.String("log-level"),
		Notify:         app.Bool("notify"),
		OutputFilename: "diffLog2Email.html",
		OutputFormat:   app.String("output-format"),
		Pass:           app.String("pass"),
		Pass2:          app.String("pass2"),
		Quite:          app.Bool("quite"),
		Reverse:        app.Bool("reverse"),
		SampleEntries:  app.Int("entries"),
		Sort:           app.String("sort"),
		SortbyCol:      app.Int("sortby-col"),
	}

	d, err := NewDB(opts)
	if err != nil {
		return nil, err
	}
	d.InitGetLogger(opts.LogLevel)
	if d.Options.LogLevel == "debug" {
		fmt.Printf("opts: \n%v\n", opts.String())
	}

	// d.setupDefaultKeyKdbx()

	// Note: credsFile used by cmds: [ add, create ]
	// fix it only this is mentioned to preserve the default value as is for defualt case
	if d.Options.Key != "" {
		credsFile = strings.Split(filepath.Base(d.Options.Key), ".")[0] + ".creds"
		credsFile = filepath.Join(filepath.Dir(d.Options.Key), credsFile)
	}

	return d, nil
}

func runCreate(app *cli.Context) error {

	d, err := newObject(app)
	if err != nil {
		fmt.Printf("failed to create db : %v\n", err)
		return err
	}

	if d.Options.Database == "" {
		d.Options.Database = "./tmp/master-db.kdbx"
	}
	if d.Options.Pass == "" {
		d.Options.Pass = gofakeit.Password(true, true, true, true, false, 16)
	}
	if d.Options.Key == "" {
		d.Options.Key = "./tmp/master-db.key"
	}

	if d.Options.SampleEntries == 0 {
		d.Options.SampleEntries = rand.Intn(12)
	}

	err = d.PreVerifyCreate()
	if err != nil {
		return err
	}

	return d.CreateKDBX()
}

func runDiff(app *cli.Context) error {

	opts := Options{
		Pass:           app.String("pass"),
		Pass2:          app.String("pass2"),
		Database:       app.String("database"),
		Database2:      app.String("database2"),
		BackupDIR:      app.String("backup-dir"),
		Key:            app.String("keyfile"),
		Key2:           app.String("keyfile2"),
		Notify:         app.Bool("notify"),
		OutputFormat:   "csv",
		OutputFilename: "diffLog2Email.html",
	}

	pattern := strings.Split(path.Base(opts.Database), ".")[0]
	if opts.Database2 == "" {
		opts.Database2 = getRecentFile(opts.BackupDIR, pattern)
	}

	diff := NewDiff(opts)
	return diff.Diff()
}

func runLs(app *cli.Context) error {

	d, err := newObject(app)
	if err != nil {
		fmt.Printf("failed to create db : %v\n", err)
		return err
	}
	if err = d.Unlock(); err != nil {
		fmt.Printf("failed to unlock dbfile: %v, err: %v\n", d.Options.Database, err)
		return err
	}
	d.FetchDBEntries()

	return d.List()
}
