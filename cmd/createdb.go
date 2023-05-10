package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tobischo/gokeepasslib/v3"
)

// func ReadPassword(item string) string {
// 	password := "super_secret"
// 	return password
// }

// NewDB creates and returns a new kdbx database
func (d *db) Create() error {

	// return if database already exist
	if _, err := os.Stat(d.Options.Database); !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%s %s already exist, won't OVERWRITE%s",
			colorRed, d.Options.Database, colorReset)
	}

	// return if keyfile already exist
	if _, err := os.Stat(d.Options.Key); !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%sfile: %s already exist, won't OVERWRITE\n%s",
			colorRed, d.Options.Key, colorReset)
	}

	// write the password to file: {database}.creds
	passFile := strings.Split(filepath.Base(d.Options.Key), ".")[0] + ".creds"
	if _, err := os.Stat(passFile); !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("%sWill overwrite password file: %s\n%s",
			colorGreen, passFile, colorReset)
	}

	// create db with some sample entries
	d.CreateWithSampleEntries()

	return nil
}

func (d *db) CreateWithSampleEntries() error {
	// fmt.Printf("here... %d, opts: %#v\n", 1, d.Options)

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

	for i := 0; i < d.Options.SampleEntries; i++ {
		// create an entry and append to rootGroup
		entry1 := CreateNewEntry("", "", "")
		rootGroup.Entries = append(rootGroup.Entries, entry1)
	}

	// create a subgroup
	subGroup := gokeepasslib.NewGroup()
	subGroup.Name = "sub group"

	for i := 0; i < d.Options.SampleEntries; i++ {
		// create an entry and append to subGroup
		entry2 := CreateNewEntry("", "", "")
		subGroup.Entries = append(subGroup.Entries, entry2)
	}

	// add subgroups to root group
	rootGroup.Groups = append(rootGroup.Groups, subGroup)

	// write keyfile and password file
	os.WriteFile(d.Options.Key, []byte("my keyfile content"), 0600)

	credsFile := strings.Split(filepath.Base(d.Options.Key), ".")[0] + ".creds"
	content := "export KDBX_DATABASE=" + d.Options.Database + "\n"
	content += "export KDBX_PASSWORD='" + d.Options.Pass + "'\n"
	content += "export KDBX_KEYFILE=" + d.Options.Key + "\n"
	os.WriteFile(credsFile, []byte(content), 0600)

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

	log.Printf("Wrote %d entries to kdbx file: %s", d.Options.SampleEntries, d.Options.Database)
	return nil
}
