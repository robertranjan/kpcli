package cmd

import (
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/sirupsen/logrus"
	"github.com/tobischo/gokeepasslib/v3"
	w "github.com/tobischo/gokeepasslib/v3/wrappers"
)

// var logger log.Logger

// Note: credsFile used by cmds: [ add, create ]
var credsFile string = "./tmp/master-db.creds"

var (
	colorGreen  = "\033[32m"
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"

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
// func (d *db) SetupLogger() {

// 	switch d.Options.LogLevel {
// 	case "debug":
// 		logger = level.NewFilter(
// 			log.NewJSONLogger(log.NewSyncWriter(os.Stdout)), level.AllowDebug(),
// 		)
// 	case "info":
// 		logger = level.NewFilter(
// 			log.NewJSONLogger(log.NewSyncWriter(os.Stdout)), level.AllowInfo(),
// 		)
// 	case "warn":
// 		logger = level.NewFilter(
// 			log.NewJSONLogger(log.NewSyncWriter(os.Stdout)), level.AllowWarn(),
// 		)
// 	default:
// 		logger = level.NewFilter(
// 			log.NewJSONLogger(log.NewSyncWriter(os.Stdout)), level.AllowError(),
// 		)
// 	}
// 	// Log some messages at different levels
// 	// level.Debug(logger).Log("message", "This is a debug message", "value", 123)
// 	// level.Info(logger).Log("message", "This is an info message", "value", 456)
// 	// level.Warn(logger).Log("message", "This is a warning message", "value", 789)
// 	// level.Error(logger).Log("message", "This is an error message", "value", 999)

// }

var log *logrus.Logger

func (d *db) InitGetLogger(logLvl string) *logrus.Logger {
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
		// CallerPrettyfier: func(f *runtime.Frame) (string, string) {
		// 	s := strings.Split(f.Function, ".")
		// 	funcname := fmt.Sprintf("%s: ", s[len(s)-1])
		// 	dir, filename := path.Split(f.File)
		// 	filename = fmt.Sprintf(" %s:%d, ", filepath.Join(filepath.Base(dir), filename), f.Line)
		// 	rs := fmt.Sprintf("%-50s", filename+funcname)
		// 	return rs, ""
		// },
	}
	if defaultLvl != log.GetLevel() {
		log.Info("changing loglevel from ", defaultLvl, " to: ", log.GetLevel())
	}
	return log
}
