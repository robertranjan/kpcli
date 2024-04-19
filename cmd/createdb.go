package cmd

import (
	"errors"
	"fmt"

	// "log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitfield/script"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/tobischo/gokeepasslib/v3"
)

// NewDB creates and returns a new kdbx database
func (d *db) PreVerifyCreate() error {

	// return if database already exist
	if _, err := os.Stat(d.Options.Database); !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%s %s already exist, won't OVERWRITE%s",
			colorRed, d.Options.Database, colorReset)
	}

	// return if keyfile already exist
	if !d.Options.NoKey {
		if _, err := os.Stat(d.Options.Key); !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("%sfile: %s already exist, won't OVERWRITE\n%s",
				colorRed, d.Options.Key, colorReset)
		}
	}

	// write the password to file: {database}.creds
	passFile := strings.Split(filepath.Base(d.Options.Database), ".")[0] + ".creds"
	if _, err := os.Stat(passFile); !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("%sWill overwrite password file: %s\n%s",
			colorGreen, passFile, colorReset)
	}

	// create db with some sample entries

	return nil
}

func GenerateKDBXEntries(n int) []gokeepasslib.Entry {
	var rv []gokeepasslib.Entry
	for i := 0; i < n; i++ {
		rv = append(rv, CreateNewEntry("", "", ""))
	}
	return rv
}

func (d *db) generateSampleEntries() gokeepasslib.Group {
	// create root group
	rootGroup := gokeepasslib.NewGroup()
	rootGroup.Name = "root group"

	rootGroup.Entries = append(rootGroup.Entries, GenerateKDBXEntries(d.Options.SampleEntries)...)

	// create a subgroup
	subGroup := gokeepasslib.NewGroup()
	subGroup.Name = "sub group"

	subGroup.Entries = append(subGroup.Entries, GenerateKDBXEntries(d.Options.SampleEntries)...)

	// add subgroups to root group
	rootGroup.Groups = append(rootGroup.Groups, subGroup)
	return rootGroup
}

func (d *db) writeCredentialsFile() error {

	content := "export KDBX_DATABASE=" + d.Options.Database + "\n"
	content += "export KDBX_PASSWORD='" + d.Options.Pass + "'\n"
	if !d.Options.NoKey {
		content += "export KDBX_KEYFILE=" + d.Options.Key + "\n"
	}
	if err := os.WriteFile(credsFile, []byte(content), 0600); err != nil {
		return fmt.Errorf("failed to write file %v, err: %v", credsFile, err)
	}
	return nil
}

func (d *db) generateCredentials() error {

	var cred *gokeepasslib.DBCredentials
	var err error

	if d.Options.NoKey {
		cred = gokeepasslib.NewPasswordCredentials(d.Options.Pass)
		if cred == nil {
			return fmt.Errorf("failed to create credentials with pass:%q ",
				d.Options.Pass)
		}
	} else {
		// Check if the keyfile exists
		if _, err = os.Stat(d.Options.Key); err == nil {
			cred, err = gokeepasslib.NewPasswordAndKeyCredentials(d.Options.Pass, d.Options.Key)
			if err != nil {
				return fmt.Errorf("failed to create credentials with pass:%q and keyFile:%q, err: %v",
					d.Options.Pass, d.Options.Key, err)
			}
		} else if os.IsNotExist(err) {
			return fmt.Errorf("%v file does not exist. \nDid you forget to mention --nokey option?", d.Options.Key)
		} else {
			return fmt.Errorf("failed to check %v file exist or not:%v", d.Options.Key, err)
		}
	}
	d.Credentials = cred
	return nil
}

func (d *db) CreateKDBX() error {

	err := os.MkdirAll(filepath.Dir(d.Options.Database), 0755)
	if err != nil {
		return err
	}
	file, err := os.Create(d.Options.Database)
	if err != nil {
		return fmt.Errorf("failed to create dbfile: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close file %s, err: %v", file.Name(), err)
		}
	}()

	// write keyfile and password file
	if !d.Options.NoKey {
		err = os.WriteFile(d.Options.Key, []byte(gofakeit.BitcoinPrivateKey()), 0600)
		if err != nil {
			return fmt.Errorf("failed to write keyfile: %v, err: %v", d.Options.Key, err)
		}
	}

	if err = d.generateCredentials(); err != nil {
		return err
	}
	if err := d.writeCredentialsFile(); err != nil {
		return err
	}

	// now create the database with the sample rootGroup
	rootGroup := d.generateSampleEntries()
	db := &gokeepasslib.Database{
		Header:      gokeepasslib.NewHeader(),
		Credentials: d.Credentials,
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
		return fmt.Errorf("failed to encode db file: %v", err)
	}

	fmt.Printf(`Created %s file with %d sample entries. To list entries,
	1. source %v
	2. kpcli ls`, d.Options.Database, d.Options.SampleEntries*2, credsFile)
	return nil
}

func (d *db) AddEntry() error {
	err := d.Unlock()
	if err != nil {
		log.Print("failed to unlock db, err: ", err)
		return err
	}
	newFile := d.Options.Database

	// create an entry and add to db
	entry1 := CreateNewEntry(d.Options.EntryTitle, d.Options.EntryUser, d.Options.EntryPass)
	rootgp := d.RawData.Content.Root.Groups[0]
	rootgp.Entries = append(d.RawData.Content.Root.Groups[0].Entries, entry1)
	d.RawData.Content.Root.Groups[0] = rootgp

	// Lock entries using stream cipher
	if err := d.RawData.LockProtectedEntries(); err != nil {
		log.Printf("error in Locking protected entries, err: %v", err)
	}

	file, err := os.Create(newFile)
	if err != nil {
		return fmt.Errorf("failed open database %q file: %v", d.Options.Database, err)
	}
	// and encode it into the file
	keepassEncoder := gokeepasslib.NewEncoder(file)
	if err := keepassEncoder.Encode(d.RawData); err != nil {
		return fmt.Errorf("failed to encode db, err: %v", err)
	}

	err = CopyFile(d.Options.Database, filepath.Join(d.Options.BackupDIR, filepath.Base(d.Options.Database)))
	if err != nil {
		return err
	}
	if !d.Options.NoKey {
		err = CopyFile(d.Options.Key, filepath.Join(d.Options.BackupDIR, filepath.Base(d.Options.Key)))
		if err != nil {
			return err
		}
	}
	err = CopyFile(credsFile, filepath.Join(d.Options.BackupDIR, filepath.Base(credsFile)))
	if err != nil {
		return err
	}

	log.Printf("kdbx with added entry(%v) has written to: %s. Total entries: %v\n",
		entry1.GetTitle(), newFile, len(rootgp.Entries))
	return nil
}

func CopyFile(cur, new string) error {
	err := os.MkdirAll(filepath.Dir(new), 0755)
	if err != nil {
		return fmt.Errorf("failed to create folder: %v", err)
	}
	cmd := "cp " + cur + " " + new
	_, err = script.Exec(cmd).Stdout()
	if err != nil {
		return fmt.Errorf("failed at script.Exec: %v", err)
	}
	log.Println("copy cmd: ", cmd)
	return nil
}
