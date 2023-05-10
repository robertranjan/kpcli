package cmd

import (
	"fmt"
	"log"
	"math/rand"
	"path"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/urfave/cli/v2"
)

func runAddEntry(app *cli.Context) error {

	// copts := Options{
	// 	Reverse:      app.Bool("reverse"),
	// 	Days:         app.Int("days"),
	// 	Pass:         app.String("pass"),
	// 	Database:     app.String("database"),
	// 	Key:          app.String("keyfile"),
	// 	Sort:         app.String("sort"),
	// 	Fields:       app.String("fields"),
	// 	SortbyCol:    app.Int("sortby-col"),
	// 	CacheFile:    app.String("cachefile"),
	// 	Quite:        app.Bool("quite"),
	// 	OutputFormat: app.String("output-format"),
	// }
	opts := Options{
		EntryTitle: app.String("title"),
		EntryUser:  app.String("user"),
		EntryPass:  app.String("entry-pass"),
		Database:   app.String("database"),
		Key:        app.String("key"),
	}

	if opts.Pass == "" {
		log.Fatalf("This command is using:\n"+
			"   database: %v\n"+
			"   keyfile: %v\n"+
			"   \033[33mCould not find password for database.\033[0m\n"+
			"   Use right cli options or export necessary env vars and try\n",
			opts.Database, opts.Key)
	}

	db, err := NewDB(opts)
	if err != nil {
		return err
	}

	return db.AddEntry()
}

func localCreateDB(app *cli.Context) (*db, error) {

	opts := Options{
		CacheFile:      app.String("cachefile"),
		Database:       app.String("database"),
		Database2:      app.String("database2"),
		Days:           app.Int("days"),
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
	d.SetupLogger()

	return d, nil
}

func runCreate(app *cli.Context) error {

	d, err := localCreateDB(app)
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

	// level.Debug(logger).Log("------> opts: ", d.Options)

	// err = d.Create()
	// if err != nil {
	// 	level.Debug(logger).Log("---->failed to createdb, err: %v", err)
	// }

	return d.Create()
}

func runDiff(app *cli.Context) error {

	opts := Options{
		Pass:           app.String("pass"),
		Pass2:          app.String("pass2"),
		Database:       app.String("database"),
		Database2:      app.String("database2"),
		Key:            app.String("keyfile"),
		Key2:           app.String("keyfile2"),
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
			"   User right cli options or export them and try again\033[0m\n",
			opts.Database, opts.Database2, opts.Key)
	}
	diff := NewDiff(opts)
	return diff.Diff()
}

func runLs(app *cli.Context) error {

	d, err := localCreateDB(app)
	if err != nil {
		fmt.Printf("failed to create db : %v\n", err)
		return err
	}
	d.Unlock()
	d.FetchDBEntries()

	return d.List()
}
