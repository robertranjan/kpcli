package createdb

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tobischo/gokeepasslib/v3"
	w "github.com/tobischo/gokeepasslib/v3/wrappers"
)

var (
	colorGreen = "\033[32m"
	colorRed   = "\033[31m"
	colorReset = "\033[0m"
	// colorYellow = "\033[33m"
)

func mkValue(key string, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{
		Key: key,
		Value: gokeepasslib.V{
			Content:   value,
			Protected: w.NewBoolWrapper(false),
		},
	}
}
func mkProtectedValue(key string, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{
		Key: key,
		Value: gokeepasslib.V{
			Content:   value,
			Protected: w.NewBoolWrapper(true),
		},
	}
}

func ReadPassword(item string) string {
	password := "super_secret"
	return password
}

func (d *db) createWithSampleEntries() error {

	err := os.MkdirAll(filepath.Dir(d.Options.Database), 0755)
	if err != nil {
		return err
	}
	file, err := os.Create(d.Options.Database)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close file %s, err: %v", file.Name(), err)
		}
	}()

	// create root group
	rootGroup := gokeepasslib.NewGroup()
	rootGroup.Name = "root group"

	entry := gokeepasslib.NewEntry()
	entry.Values = append(entry.Values, mkValue("Title", "My GMail password"))
	entry.Values = append(entry.Values, mkValue("UserName", "example@gmail.com"))
	entry.Values = append(entry.Values, mkProtectedValue("Password", "hunter2"))

	rootGroup.Entries = append(rootGroup.Entries, entry)

	// demonstrate creating sub group (we'll leave it empty because we're lazy)
	subGroup := gokeepasslib.NewGroup()
	subGroup.Name = "sub group"

	subEntry := gokeepasslib.NewEntry()
	subEntry.Values = append(subEntry.Values, mkValue("Title", "Another password"))
	subEntry.Values = append(subEntry.Values, mkValue("UserName", "johndough"))
	subEntry.Values = append(subEntry.Values, mkProtectedValue("Password", "123456"))

	subGroup.Entries = append(subGroup.Entries, subEntry)

	rootGroup.Groups = append(rootGroup.Groups, subGroup)

	// write keyfile and password file
	os.WriteFile(d.Options.Key, []byte("my keyfile content"), 0600)
	passFile := strings.Split(filepath.Base(d.Options.Key), ".")[0] + ".password"
	os.WriteFile(passFile, []byte(d.Options.Pass), 0600)

	cred, err := gokeepasslib.NewPasswordAndKeyCredentials(d.Options.Pass, d.Options.Key)
	if err != nil {
		return fmt.Errorf("failed to create credentials, err: %v", err)
	}
	// now create the database containing the root group
	db := &gokeepasslib.Database{
		Header:      gokeepasslib.NewHeader(),
		Credentials: cred,
		Content: &gokeepasslib.DBContent{
			Meta: gokeepasslib.NewMetaData(),
			Root: &gokeepasslib.RootData{
				Groups: []gokeepasslib.Group{rootGroup},
			},
		},
	}

	// Lock entries using stream cipher
	if err := db.LockProtectedEntries(); err != nil {
		log.Printf("error in Locking protected entries, err: %v", err)
	}

	// and encode it into the file
	keepassEncoder := gokeepasslib.NewEncoder(file)
	if err := keepassEncoder.Encode(db); err != nil {
		panic(err)
	}

	log.Printf("Wrote kdbx file: %s", d.Options.Database)
	return nil
}

// NewDB creates and returns a new kdbx database
func NewDB(opt Options) (*db, error) {
	d := &db{
		Options: &opt,
	}
	// return if database already exist
	if _, err := os.Stat(d.Options.Database); !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("%s %s already exist, won't OVERWRITE%s",
			colorRed, d.Options.Database, colorReset)
	}

	// return if keyfile already exist
	if _, err := os.Stat(d.Options.Key); !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("%sfile: %s already exist, won't OVERWRITE\n%s",
			colorRed, d.Options.Key, colorReset)
	}
	// create db with some sample entries
	d.createWithSampleEntries()

	// write the password to file: {database}.password
	passFile := strings.Split(filepath.Base(d.Options.Key), ".")[0] + ".password"
	if _, err := os.Stat(passFile); !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("%sWill overwrite password file: %s\n%s",
			colorGreen, passFile, colorReset)
	}
	return d, nil
}
