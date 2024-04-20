package models

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tobischo/gokeepasslib/v3"
)

type Client struct {
	Options        *Options
	Credentials    *gokeepasslib.DBCredentials
	CredentialFile string
}

// Options holds the cli options
type Options struct {
	//diff
	BackupDIR      string
	CacheFile      string // ls
	Database       string
	Database2      string
	DiffCalling    bool
	Notify         bool
	OutputFilename string
	Pass           string
	Pass2          string

	// add entry
	EntryPass  string
	EntryTitle string
	EntryUser  string

	// ls
	Days      int
	Fields    string
	Key       string
	Key2      string
	Reverse   bool
	Sort      string
	SortbyCol int

	// common
	Config       string
	LogLevel     string
	OutputFormat string // ls,diff
	Quite        bool   // ls

	// create db
	NoKey         bool
	SampleEntries int
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

type Diff struct {
	ToDBOption     *Options
	FromDBOption   *Options
	Options        *Options
	OutputFilename string
}

func (o *Options) String() string {
	d, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		log.Debugf("failed to marshal option, err: %v", err)
	}
	return string(d)
}
