package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/tobischo/gokeepasslib/v3"
)

type Config struct {
	Notify
	Create
	DiffCfg
}

func (c *Config) String() string {
	s, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "failed to marshal config"
	}
	return string(s)
}

func (c *Config) loadFromFile(filename string) error {

	// Check if the keyfile exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("%v file does not exist", filename)
	}

	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	if err = toml.Unmarshal(b, &config); err != nil {
		return fmt.Errorf("unmarshall failed")
	}

	return nil
}

type DiffCfg struct {
	Database1      string `toml:"database1"`
	Database2      string `toml:"database2"`
	Keyfile1       string `toml:"keyfile1"`
	Keyfile2       string `toml:"keyfile2"`
	OutputFilename string `toml:"outputFilename"`
	Password1      string `toml:"password1"`
	Password2      string `toml:"password2"`
}

type Notify struct {
	EmailContent  string   `toml:"emailContent"`
	From          string   `toml:"from"`
	EmailPassword string   `toml:"emailPassword"`
	SMTPHost      string   `toml:"smtpHost"`
	SMTPPort      int      `toml:"smtpPort"`
	Subject       string   `toml:"subject"`
	To            []string `toml:"to"`
}

type Create struct {
	Databaese string `toml:"databaese"`
	Keyfile   string `toml:"keyfile"`
	Password  string `toml:"password"`
}

type Diff struct {
	ToDBOption     *Options
	FromDBOption   *Options
	options        *Options
	OutputFilename string
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
}

// Options holds the cli options
type Options struct {
	BackupDIR      string
	CacheFile      string
	Config         string
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
