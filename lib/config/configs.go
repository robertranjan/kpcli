package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/robertranjan/kpcli/lib/utils"
)

type Config struct {
	Notify
	Create
	DiffCfg
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

func New(filename string) (*Config, error) {
	// Check if the keyfile exists
	if utils.IsFileNotExist(filename) {
		return nil, fmt.Errorf("%v file does not exist", filename)
	}

	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var c Config
	if err = toml.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("unmarshall failed")
	}
	return &c, nil
}

func (c *Config) String() string {
	s, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "failed to marshal config"
	}
	return string(s)
}

var SampleConfig = `
[notify]
emailContent = "will be generated during execution"
emailPassword = "keepass_gmail_app_password"
from = "yourEmail@gmail.com"
smtpHost = "smtp.gmail.com"
smtpPort = 587
subject = "here are the KDBX changes since last backup!"
to = ["yourEmail@gmail.com", "email2@domain.com"]

[create]
database = "./tmp/master-db.kdbx"
keyfile = "./tmp/master-db.key"
password = "super_s3cr3t"

[diffCfg]
database1 = "./tmp/database1"
database2 = "./tmp/database2"
keyfile1 = "./tmp/keyfile1"
keyfile2 = "./tmp/keyfile2"
outputFilename = "diffLog2Email.html"
password1 = "super_secret"
password2 = "super_secret"
`
