package cmd

import (
	"log"
	"os"

	"github.com/tobischo/gokeepasslib/v3"
)

func (d *db) AddEntry() error {
	d.Unlock()
	newFile := d.Options.Database + ".tmp.kdbx"

	// create an entry and append to subGroup
	entry2 := CreateNewEntry("", "", "")

	rootgp := d.RawData.Content.Root.Groups[0]
	rootgp.Entries = append(rootgp.Entries, entry2)

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

	log.Printf("Wrote kdbx file: %s", newFile)
	return nil
}
