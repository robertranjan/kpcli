package cmd

import (
	"fmt"
	"strconv"
	"time"

	// "log"
	"os"
	"path/filepath"
	"strings"

	"github.com/robertranjan/kpcli/lib/models"

	"github.com/bitfield/script"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/robertranjan/kpcli/lib/utils"
	"github.com/tobischo/gokeepasslib/v3"
)

type kdbx struct {
	Entries         []models.Interested
	Options         *models.Options
	SelectedEntries []models.Interested
	RawData         *gokeepasslib.Database
}

var db *kdbx

// NewDB creates and returns a new kdbx database
func (c *Client) PreVerifyCreate() error {
	// return if database already exist
	if utils.IsFileExist(c.Options.Database) {
		return fmt.Errorf("%s %s already exist, won't OVERWRITE%s",
			utils.ColorRed, c.Options.Database, utils.ColorReset)
	}

	// return if keyfile already exist
	if !c.Options.NoKey {
		if utils.IsFileExist(c.Options.Key) {
			return fmt.Errorf("%sfile: %s already exist, won't OVERWRITE\n%s",
				utils.ColorRed, c.Options.Key, utils.ColorReset)
		}
	}

	// write the credentials to file: {database}.creds
	cred := strings.Split(filepath.Base(c.Options.Database), ".")[0] + ".creds"
	if utils.IsFileExist(cred) {
		fmt.Printf("%sWill overwrite password file: %s\n%s",
			utils.ColorGreen, cred, utils.ColorReset)
	}
	return nil
}

func GenerateKDBXEntries(n int) []gokeepasslib.Entry {
	var rv []gokeepasslib.Entry
	for i := 0; i < n; i++ {
		rv = append(rv, CreateNewEntry("", "", ""))
	}
	return rv
}

func (d *kdbx) generateSampleEntries() gokeepasslib.Group {
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

func (c *Client) writeCredentialsFile() error {
	content := "export KDBX_DATABASE=" + c.Options.Database + "\n"
	content += "export KDBX_PASSWORD='" + c.Options.Pass + "'\n"
	if !c.Options.NoKey {
		content += "export KDBX_KEYFILE=" + c.Options.Key + "\n"
	}
	if err := os.WriteFile(c.CredentialFile, []byte(content), 0600); err != nil {
		return fmt.Errorf("failed to write file %v, err: %v", c.CredentialFile, err)
	}
	return nil
}

func (c *Client) generateCredentials() error {
	var cred *gokeepasslib.DBCredentials
	var err error

	if c.Options.NoKey {
		cred = gokeepasslib.NewPasswordCredentials(c.Options.Pass)
		if cred == nil {
			return fmt.Errorf("failed to create credentials with pass:%q ",
				c.Options.Pass)
		}
		c.Credentials = cred
		return nil
	}
	// check keyfile
	if !utils.IsFileExist(c.Options.Key) {
		return fmt.Errorf("%v file does not exist. \nDid you forget to mention --nokey option?", c.Options.Key)
	}
	// gen creds with keyfile
	cred, err = gokeepasslib.NewPasswordAndKeyCredentials(c.Options.Pass, c.Options.Key)
	if err != nil {
		return fmt.Errorf("failed to create credentials with pass:%q and keyFile:%q, err: %v",
			c.Options.Pass, c.Options.Key, err)
	}

	c.Credentials = cred
	return nil
}

func (c *Client) CreateKDBX() (*gokeepasslib.Database, error) {
	err := os.MkdirAll(filepath.Dir(c.Options.Database), 0755)
	if err != nil {
		return nil, err
	}
	file, err := os.Create(c.Options.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to create dbfile: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close file %s, err: %v", file.Name(), err)
		}
	}()

	// write keyfile and password file
	if !c.Options.NoKey {
		err = os.WriteFile(c.Options.Key, []byte(gofakeit.BitcoinPrivateKey()), 0600)
		if err != nil {
			return nil, fmt.Errorf("failed to write keyfile: %v, err: %v", c.Options.Key, err)
		}
	}

	if err = c.generateCredentials(); err != nil {
		return nil, err
	}
	if err := c.writeCredentialsFile(); err != nil {
		return nil, err
	}

	// now create the database with the sample rootGroup
	rootGroup := db.generateSampleEntries()
	database := &gokeepasslib.Database{
		Header:      gokeepasslib.NewHeader(),
		Credentials: c.Credentials,
		Content: &gokeepasslib.DBContent{
			Meta: gokeepasslib.NewMetaData(),
			Root: &gokeepasslib.RootData{
				Groups: []gokeepasslib.Group{rootGroup},
			},
		},
	}

	// Lock entries using stream cipher
	if err := database.LockProtectedEntries(); err != nil {
		log.Printf("error in Locking protected entries, err: %v", err)
	}

	// and encode it into the file
	keepassEncoder := gokeepasslib.NewEncoder(file)
	if err := keepassEncoder.Encode(database); err != nil {
		return nil, fmt.Errorf("failed to encode db file: %v", err)
	}

	fmt.Printf(`Created %s file with %d sample entries. To list entries,
	1. source %v
	2. kpcli ls`, c.Options.Database, c.Options.SampleEntries*2, c.CredentialFile)
	return database, nil
}

func (d *kdbx) AddEntry() error {
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
		return fmt.Errorf("failed to create database %q file: %v", d.Options.Database, err)
	}
	// and encode it into the file
	keepassEncoder := gokeepasslib.NewEncoder(file)
	if err := keepassEncoder.Encode(d.RawData); err != nil {
		return fmt.Errorf("failed to encode db, err: %v", err)
	}

	// make a copy/backup of kdbx database to backupDir
	err = backupFile(d.Options.Database)
	if err != nil {
		return err
	}

	// make a copy/backup of keyfile to backupDir
	if !d.Options.NoKey {
		err = backupFile(d.Options.Key)
		if err != nil {
			return err
		}
	}

	log.Debugf("kdbx with added entry(%v) has written to: %s. Total entries: %v\n",
		entry1.GetTitle(), newFile, len(rootgp.Entries))
	return nil
}

func backupFile(cur string) error {
	d, f := filepath.Split(cur)
	newFile := filepath.Join(d, f+"."+strconv.Itoa(time.Now().Nanosecond()))

	cmd := "cp " + cur + " " + newFile
	_, err := script.Exec(cmd).Stdout()
	if err != nil {
		return fmt.Errorf("failed at script.Exec: %v", err)
	}
	log.Debug("copy cmd: ", cmd)
	return nil
}
