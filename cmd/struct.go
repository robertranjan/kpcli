package cmd

import (
	"fmt"
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
	// V               *viper.Viper
}

func (o *Options) String() string {
	return fmt.Sprintf("CacheFile: %v, Database: %v, Days: %v, Diff: %v, "+
		"Fields: %v, Key: %v, OutputFormat: %v, Pass: %v, Quite: %v, "+
		"Reverse: %v, Sort: %v, SortbyCol: %v, Title: %v, User: %v, Pass: %v",
		o.CacheFile, o.Database, o.Days,
		o.DiffCalling, o.Fields, o.Key, o.OutputFormat, "****",
		o.Quite, o.Reverse, o.Sort, o.SortbyCol,
		o.EntryTitle, o.EntryUser, o.EntryPass,
	)
}

// Options holds the cli options
type Options struct {
	CacheFile      string
	Database       string
	Database2      string
	Days           int
	LogLevel       string
	DiffCalling    bool
	EntryPass      string
	EntryTitle     string
	EntryUser      string
	Fields         string
	Key            string
	Key2           string
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
