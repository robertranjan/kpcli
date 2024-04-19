package cmd

import (
	"fmt"
	"math/rand"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/sirupsen/logrus"
	"github.com/tobischo/gokeepasslib/v3"
	w "github.com/tobischo/gokeepasslib/v3/wrappers"
	"github.com/urfave/cli/v2"
)

var (
	// Note: credsFile used by cmds: [ add, create ]
	credsFile   string = "./tmp/master-db.creds"
	colorGreen         = "\033[32m"
	colorReset         = "\033[0m"
	colorRed           = "\033[31m"
	colorYellow        = "\033[33m"
	TimeLayout         = "2006-01-02 15:04:05"
	lengthUser         = 25
	configFile         = "unavailable_kpcli.toml"
	config      Config
	log         *logrus.Logger
)

var CmdAdd = &cli.Command{
	Name:    "add",
	Usage:   "add an entry",
	Aliases: []string{"a"},
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
		create
		
	kpcli \
		--nokey \
		--pass 'super_secret' \
		--db ./tmp/master-db.kdbx \
		create
		
		`,

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
		&cli.StringFlag{
			Name:    "backup-dir",
			Usage:   "dir to look for recent backup file(when database2 is not given)",
			Value:   "./bkups/",
			Aliases: []string{"bkup"},
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
		ls

	kpcli \
		--nokey \
		--pass 'super_secret' \
		--db ./tmp/master-db.kdbx \
		ls

	kpcli \
		--keyfile ./tmp/master-db.key \
		--pass 'super_secret' \
		--db ./tmp/master-db.kdbx \
		ls \
		--fields few

	Fields options:
		all		: {cols[0], cols[1], cols[2], cols[3], cols[4]}
		few		: {cols[0]}
		default	: {cols[0], cols[1], cols[2], cols[3]}
		`,

	Action: runLs,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "fields",
			Usage:   "fields list to be displayed, available options: all|few",
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

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

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

func loadConfigOnDemand(configFile string) {
	if configFile == "" {
		// no need to load as configFile == null
		return
	}
	if err := config.loadFromFile(configFile); err != nil {
		fmt.Println("LoadConfig failed. Continuing with cli args...")
	}
	log.Debugf("config: %s\n", config.String())
}

func getOutputFilename() string {
	outputFilename := "diffLog2Email.html"
	if config.OutputFilename != "" {
		outputFilename = config.OutputFilename
	}
	return outputFilename
}

func newObject(app *cli.Context) (*db, error) {

	//setup logger
	log = InitGetLogger(app.String("log-level"))

	// read config on demand
	configFile = app.String("config")
	loadConfigOnDemand(configFile)
	outputFilename := getOutputFilename()

	opts := Options{
		CacheFile:      app.String("cachefile"),
		Config:         app.String("config"),
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
		NoKey:          app.Bool("nokey"),
		Notify:         app.Bool("notify"),
		OutputFilename: outputFilename,
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

	if d.Options.LogLevel == "debug" {
		log.Printf("opts: \n%v\n", opts.String())
	}

	// generate credsFile path
	if d.Options.Key != "" {
		credsFile = strings.Split(filepath.Base(d.Options.Key), ".")[0] + ".creds"
		credsFile = filepath.Join(filepath.Dir(d.Options.Key), credsFile)
	}
	if d.Options.NoKey {
		credsFile = strings.Split(filepath.Base(d.Options.Database), ".")[0] + ".creds"
		credsFile = filepath.Join(filepath.Dir(d.Options.Database), credsFile)
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
		d.Options.SampleEntries = rand.Intn(11) + 1
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

	if d.Options.Key == "" || d.Options.Database == "" {
		fmt.Printf("%v"+`   --database is a required arguments.
	If you are trying, run below commands:
	1. kpcli createdb
	2. kpcli ls`+"%v\nHere is usage:\n%v", colorYellow, colorGreen, colorReset)
		cli.ShowAppHelpAndExit(app, 0)
		return nil
	}

	if err = d.Unlock(); err != nil {
		fmt.Printf("failed to unlock dbfile: %v, err: %v\n", d.Options.Database, err)
		return err
	}
	d.FetchDBEntries()

	return d.List()
}

// NewDB create and return a new kdbx db object
func NewDB(opts Options) (*db, error) {
	d := &db{
		Options: &opts,
	}
	return d, nil
}

func MkValue(key string, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{
		Key: key,
		Value: gokeepasslib.V{
			Content:   value,
			Protected: w.NewBoolWrapper(false),
		},
	}
}

func MkProtectedValue(key string, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{
		Key: key,
		Value: gokeepasslib.V{
			Content:   value,
			Protected: w.NewBoolWrapper(true),
		},
	}
}

func CreateNewEntry(t, u, p string) gokeepasslib.Entry {
	// set defaults if param is empty
	if t == "" {
		t = gofakeit.NewCrypto().AppName()
	}
	if u == "" {
		u = gofakeit.Username()
	}
	if p == "" {
		p = gofakeit.Password(true, true, true, true, false, 16)
	}
	entry := gokeepasslib.NewEntry()
	entry.Values = append(entry.Values, MkValue("Title", t))
	entry.Values = append(entry.Values, MkValue("UserName", u))
	entry.Values = append(entry.Values, MkProtectedValue("Password", p))
	return entry
}

func InitGetLogger(logLvl string) *logrus.Logger {
	log = logrus.New()
	defaultLvl := log.GetLevel()

	switch logLvl {
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "trace":
		log.SetLevel(logrus.TraceLevel)
	}

	log.SetReportCaller(true)
	log.Formatter = &logrus.TextFormatter{
		DisableTimestamp: false,
		DisableColors:    true,
	}
	if defaultLvl != log.GetLevel() {
		log.Info("changing loglevel from ", defaultLvl, " to: ", log.GetLevel())
	}
	return log
}
