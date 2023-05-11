package cmd

import (
	"math/rand"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/tobischo/gokeepasslib/v3"
	w "github.com/tobischo/gokeepasslib/v3/wrappers"
)

var logger log.Logger
var credsFile string

// BackupDIR is where the backup databases are
const BackupDIR = "./bkups/"

var (
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
	colorRed   = "\033[31m"

	TimeLayout = "2006-01-02 15:04:05"
	lengthUser = 25
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// NewDB create and return a new kdbx db object
func NewDB(opts Options) (*db, error) {
	d := &db{
		Options: &opts,
	}
	// d.Unlock()
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

// SetupLogger create and setup logger with zap pkg
func (d *db) SetupLogger() {

	switch d.Options.LogLevel {
	case "debug":
		logger = level.NewFilter(
			log.NewJSONLogger(log.NewSyncWriter(os.Stdout)), level.AllowDebug(),
		)
	case "info":
		logger = level.NewFilter(
			log.NewJSONLogger(log.NewSyncWriter(os.Stdout)), level.AllowInfo(),
		)
	case "warn":
		logger = level.NewFilter(
			log.NewJSONLogger(log.NewSyncWriter(os.Stdout)), level.AllowWarn(),
		)
	default:
		logger = level.NewFilter(
			log.NewJSONLogger(log.NewSyncWriter(os.Stdout)), level.AllowError(),
		)
	}
	// Log some messages at different levels
	level.Debug(logger).Log("message", "This is a debug message", "value", 123)
	level.Info(logger).Log("message", "This is an info message", "value", 456)
	level.Warn(logger).Log("message", "This is a warning message", "value", 789)
	level.Error(logger).Log("message", "This is an error message", "value", 999)

}
