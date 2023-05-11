package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/bitfield/script"
	"github.com/go-kit/log/level"
	"github.com/tobischo/gokeepasslib/v3"
)

func (d *db) AddEntry() error {
	err := d.Unlock()
	if err != nil {
		level.Error(logger).Log("failed to unlock db, err: ", err)
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
		log.Fatalf("failed open database file: %v, err: %v", d.Options.Database, err)
	}
	// and encode it into the file
	keepassEncoder := gokeepasslib.NewEncoder(file)
	if err := keepassEncoder.Encode(d.RawData); err != nil {
		panic(err)
	}

	CopyFile(d.Options.Database, filepath.Join(BackupDIR, filepath.Base(d.Options.Database)))
	CopyFile(d.Options.Key, filepath.Join(BackupDIR, filepath.Base(d.Options.Key)))
	CopyFile(credsFile, filepath.Join(BackupDIR, filepath.Base(credsFile)))

	log.Printf("kdbx with added entry(%v) has written to: %s. Total entries: %v\n",
		entry1.GetTitle(), newFile, len(rootgp.Entries))
	return nil
}

func CopyFile(cur, new string) {
	os.MkdirAll(filepath.Dir(new), 0755)
	cmd := "cp " + cur + " " + new
	script.Exec(cmd).Stdout()
	level.Debug(logger).Log("copy cmd: ", cmd)
}
