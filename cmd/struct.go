package cmd

import (
	"encoding/json"
	"time"

	"github.com/tobischo/gokeepasslib/v3"
)

// var sugar *zap.SugaredLogger

type Diff struct {
	ToDBOption   *Options
	FromDBOption *Options
	options      *Options
}

type db struct {
	Entries         []Interested
	Options         *Options
	SelectedEntries []Interested
	RawData         *gokeepasslib.Database
	Credentials     *gokeepasslib.DBCredentials
	// V               *viper.Viper
}

func (o *Options) String() string {

	d, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		log.Debugf("failed to marshal option, err: %v", err)
	}
	return string(d)

	// return fmt.Sprintf("CacheFile: %v, Database: %v, Days: %v, Diff: %v, "+
	// 	"Fields: %v, Key: %v, OutputFormat: %v, Pass: %v, Quite: %v, "+
	// 	"Reverse: %v, Sort: %v, SortbyCol: %v, Title: %v, User: %v, Pass: %v",
	// 	o.CacheFile, o.Database, o.Days,
	// 	o.DiffCalling, o.Fields, o.Key, o.OutputFormat, "****",
	// 	o.Quite, o.Reverse, o.Sort, o.SortbyCol,
	// 	o.EntryTitle, o.EntryUser, o.EntryPass,
	// )
}

// Options holds the cli options
type Options struct {
	BackupDIR      string
	CacheFile      string
	Database       string
	Database2      string
	Days           int
	DiffCalling    bool
	EntryPass      string
	EntryTitle     string
	EntryUser      string
	Fields         string
	Key            string
	Key2           string
	LogLevel       string
	NoKey          bool
	Notify         bool
	OutputFilename string
	OutputFormat   string
	Pass           string
	Pass2          string
	Quite          bool
	Reverse        bool
	SampleEntries  int
	Sort           string
	SortbyCol      int
}

type Interested struct {
	Created   time.Time
	Histories int
	KeyValues map[string]string
	Modified  time.Time
	Pass      string
	Tags      string
	Title     string
	User      string
}
