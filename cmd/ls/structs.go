package ls

import (
	"fmt"
	"time"
)

// Options holds the cli options
type Options struct {
	Database     string
	Days         int
	Fields       string
	Key          string
	OutputFormat string
	Pass         string
	Quite        bool
	Reverse      bool
	Sort         string
	SortbyCol    int
	// options for cmd: diff
	CacheFile   string
	DiffCalling bool
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

type db struct {
	Entries         []Interested
	Options         *Options
	SelectedEntries []Interested
}

func (o *Options) String() string {
	return fmt.Sprintf("CacheFile: %v, Database: %v, Days: %v, Diff: %v, "+
		"Fields: %v, Key: %v, OutputFormat: %v, Pass: %v, Quite: %v, "+
		"Reverse: %v, Sort: %v, SortbyCol: %v",
		o.CacheFile, o.Database, o.Days,
		o.DiffCalling, o.Fields, o.Key, o.OutputFormat, "****",
		o.Quite, o.Reverse, o.Sort, o.SortbyCol,
	)
}
