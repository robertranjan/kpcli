package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/robertranjan/kpcli/lib/config"
	"github.com/robertranjan/kpcli/lib/models"
	"github.com/robertranjan/kpcli/lib/utils"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/sirupsen/logrus"
	"github.com/tobischo/gokeepasslib/v3"
	w "github.com/tobischo/gokeepasslib/v3/wrappers"
	"github.com/urfave/cli/v2"
)

type Client models.Client

var client *Client

var (
	// Note: credsFile used by cmds: [ add, create ]
	// credsFile  = "./tmp/master-db.creds"
	TimeLayout = "2006-01-02 15:04:05"
	lengthUser = 25
	configFile = "unavailable_kpcli.toml"
	cfg        *config.Config
	log        *logrus.Logger
)

// CmdAdd helps user to add entry to password db
var CmdAdd = &cli.Command{
	Name:    "add",
	Usage:   "add an entry",
	Aliases: []string{"a"},
	Description: `Add an entry to a .kdbx database

syntax:
	./kpcli --keyfile <keyfile> \
			--database <database-filename> \
			--pass <pass to open database> \
		add --entry-title <title> \
			--entry-user <username> \
			--entry-pass <password>

Example:
		kpcli add \
			--entry-title new-entry-1 \
			--entry-user example-user1 \
			--entry-pass secret_13

		task add-entry
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

// CmdCreatedb creates password db
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
			Name:    "sample-entries",
			Usage:   "number of sample entries",
			Aliases: []string{"se"},
		},
	},
}

// CmdDiff runs diff between 2 passward dbs
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

// CmdGenerateSampleConfig generate a sample config
var CmdGenerateSampleConfig = &cli.Command{
	Name:    "generate-sample-config",
	Usage:   "generate sample config file: kpcli.toml",
	Aliases: []string{"gen"},
	Description: `generate sample config file: kpcli.toml

Example:
	kpcli \
		generate-sample-config
	`,
	Action: runGenerateSampleConfig,
}

func runGenerateSampleConfig(app *cli.Context) error {
	InitGetLogger(app.String("log-level"))
	if utils.IsFileExist("kpcli.toml") {
		log.Info("found file, backing up existing file to tmp/")
		backupFile("kpcli.toml")
	}
	fmt.Println("writing sample config file: kpcli.toml")
	return os.WriteFile("kpcli.toml", []byte(config.SampleConfig), 0600)
}

// CmdLs lists entries from a db
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

// LoadConfigOnDemand loads config if present
func LoadConfigOnDemand(configFile string) {
	if configFile == "" {
		// no need to load as configFile == null
		return
	}
	var err error
	if cfg, err = config.New(configFile); err != nil {
		fmt.Println("LoadConfig failed. Continuing with cli args...")
	}
	log.Debugf("config: %s\n", cfg.String())
}

// runAddEntry - update pkg->var: client
func runAddEntry(app *cli.Context) error {
	var err error
	opts := getOptions(app)

	client, err = newClient(opts)
	if err != nil {
		fmt.Printf("failed to create db : %v\n", err)
		return err
	}

	db = &kdbx{
		Options: client.Options,
	}

	if client.Options.Pass == "" {
		client.Options.Pass = "super_secret"
	}
	return db.AddEntry()
}

func getOutputFilename() string {
	outputFilename := "diffLog2Email.html"
	if cfg != nil && cfg.OutputFilename != "" {
		outputFilename = cfg.OutputFilename
	}
	return outputFilename
}

func getOptions(app *cli.Context) *models.Options {
	outputFilename := getOutputFilename()

	return &models.Options{
		// diff
		BackupDIR:      app.String("backup-dir"),
		CacheFile:      app.String("cachefile"),
		Database:       app.String("database"),
		Database2:      app.String("database2"),
		DiffCalling:    app.Bool("diff-calling"),
		Notify:         app.Bool("notify"),
		OutputFilename: outputFilename,
		Pass:           app.String("pass"),
		Pass2:          app.String("pass2"),

		// add entry
		EntryPass:  app.String("entry-pass"),
		EntryTitle: app.String("entry-title"),
		EntryUser:  app.String("entry-user"),

		// ls
		Days:      app.Int("days"),
		Fields:    app.String("fields"),
		Key:       app.String("keyfile"),
		Key2:      app.String("keyfile2"),
		Reverse:   app.Bool("reverse"),
		Sort:      app.String("sort"),
		SortbyCol: app.Int("sortby-col"),

		// common
		Config:       app.String("config"),
		LogLevel:     app.String("log-level"),
		OutputFormat: app.String("output-format"),
		Quite:        app.Bool("quite"),

		// create db
		NoKey:         app.Bool("nokey"),
		SampleEntries: app.Int("sample-entries"),
	}
}

// newClient creates a base db object using cli-args
func newClient(opts *models.Options) (*Client, error) {
	//setup logger
	log = InitGetLogger(opts.LogLevel)

	// read config on demand
	// configFile = app.String("config")
	configFile = opts.Config
	LoadConfigOnDemand(configFile)

	// opts := getOptions(app)
	// log.Debugf("opts: \n%v\n", opts.String())

	c := &Client{
		Options:        opts,
		CredentialFile: "./tmp/master-db.creds",
	}

	// generate credsFile path
	// overwrite default value with user-args
	credsFile := strings.Split(filepath.Base(c.Options.Database), ".")[0] + ".creds"
	credsFile = filepath.Join(filepath.Dir(c.Options.Database), credsFile)
	c.CredentialFile = credsFile
	return c, nil
}

func updateDefaultOptionValues() {
	if client.Options.Database == "" {
		client.Options.Database = "./tmp/master-db.kdbx"
	}
	if client.Options.Key == "" {
		client.Options.Key = "./tmp/master-db.key"
	}
	if client.Options.SampleEntries == 0 {
		client.Options.SampleEntries = rand.Intn(11) + 1
	}
}

func runCreate(app *cli.Context) error {
	var err error
	opts := getOptions(app)

	client, err = newClient(opts)
	if err != nil {
		fmt.Printf("failed to create db : %v\n", err)
		return err
	}

	updateDefaultOptionValues()

	if client.Options.Pass == "" {
		// generate a sample password with(lwr,  upr, numeric,spl, space, length )
		client.Options.Pass = gofakeit.Password(true, true, true, true, false, 16)
	}

	err = client.PreVerifyCreate()
	if err != nil {
		return err
	}
	// initialize global database: d
	db = &kdbx{
		Options: client.Options,
	}
	db.RawData, err = client.CreateKDBX()
	if err != nil {
		return err
	}

	return nil
}

func runDiff(app *cli.Context) error {
	opts := getOptions(app)
	// opts := models.Options{
	// 	Pass:           app.String("pass"),
	// 	Pass2:          app.String("pass2"),
	// 	Database:       app.String("database"),
	// 	Database2:      app.String("database2"),
	// 	BackupDIR:      app.String("backup-dir"),
	// 	Key:            app.String("keyfile"),
	// 	Key2:           app.String("keyfile2"),
	// 	Notify:         app.Bool("notify"),
	// 	OutputFormat:   "csv",
	// 	OutputFilename: "diffLog2Email.html",
	// }

	pattern := strings.Split(path.Base(opts.Database), ".")[0]
	if opts.Database2 == "" {
		opts.Database2 = getRecentFile(opts.BackupDIR, pattern)
	}

	diff := NewDiff(opts)
	return diff.Diff()
}

func runLs(app *cli.Context) error {
	var err error
	opts := getOptions(app)

	client, err = newClient(opts)
	if err != nil {
		fmt.Printf("failed to create client : %v\n", err)
		return err
	}

	// NewDB create and return a new kdbx db object

	db = &kdbx{
		Options: client.Options,
	}

	if err = db.Unlock(); err != nil {
		fmt.Printf("failed to unlock dbfile: %v, err: %v\n", client.Options.Database, err)
		return err
	}
	db.FetchDBEntries()

	return db.List()
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
