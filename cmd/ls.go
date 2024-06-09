package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/robertranjan/kpcli/lib/models"
	"github.com/robertranjan/kpcli/lib/utils"
	"github.com/tobischo/gokeepasslib/v3"
)

func (d *kdbx) Unlock() error {
	log.Debugf("options: %#v\n", d.Options)
	file, err := os.Open(d.Options.Database)
	if err != nil {
		return fmt.Errorf("failed open database %q file: %v", d.Options.Database, err)
	}

	db := gokeepasslib.NewDatabase(gokeepasslib.WithDatabaseKDBXVersion4())
	if err := client.generateCredentials(); err != nil {
		return err
	}
	db.Credentials = client.Credentials

	if err := gokeepasslib.NewDecoder(file).Decode(db); err != nil {
		log.Error("failed to decode dbfile: ", d.Options.Database, " err:", err)
		if d.Options.LogLevel == "debug" {
			fmt.Printf("opts: \n%v\n", d.Options.String())
		}
		return err
	}

	if err := db.UnlockProtectedEntries(); err != nil {
		log.Errorf("failed to unload db, err: %v", err)
		return err
	}
	d.RawData = db
	return nil
}

func (d *kdbx) FetchDBEntries() {
	for _, rootgp := range d.RawData.Content.Root.Groups {
		d.FetchGrpEntries(rootgp, "")
	}
}

func (d *kdbx) FetchGrpEntries(grp gokeepasslib.Group, parentPath string) {
	for _, e := range grp.Entries {
		// Construct the full title including the parent path
		fullTitle := parentPath + "/" + strings.TrimSpace(e.GetTitle())

		kv := make(map[string]string)
		for _, entry := range e.Values {
			kv[entry.Key] = entry.Value.Content
		}

		var hist int
		if len(e.Histories) > 0 {
			hist = len(e.Histories[0].Entries)
		}

		// Create the Interested struct with the full title
		et := models.Interested{
			Title:     fullTitle,
			User:      strings.TrimSpace(e.GetContent("UserName")),
			Pass:      e.GetPassword(),
			Tags:      e.Tags,
			Histories: hist,
			Created:   e.Times.CreationTime.Time,
			Modified:  e.Times.LastModificationTime.Time,
			KeyValues: kv,
		}
		if len(et.User) > lengthUser {
			et.User = et.User[:lengthUser] + "..."
		}
		d.Entries = append(d.Entries, et)
	}
	// Recursively fetch entries from subgroups
	for _, subgrp := range grp.Groups {
		d.FetchGrpEntries(subgrp, parentPath+"/"+grp.Name)
	}
}

func (d *kdbx) Display() {

	t := d.getTable()
	t = d.updateTableWithSelectedEntries(t)

	// Note:
	//    it sort with string sort. Try `--sort-by-col 2`` to see how it sorts.
	t.SortBy([]table.SortBy{
		{Number: d.Options.SortbyCol, Mode: table.Asc},
	})
	if d.Options.Reverse {
		t.SortBy([]table.SortBy{
			{Number: d.Options.SortbyCol, Mode: table.Dsc},
		})
	}

	fmt.Println()

	// when cacheFile is set
	//   `./kpcli diff` uses this to save output in file & then diff them
	if d.Options.CacheFile != "" {
		cacheFile(t, d.Options.CacheFile)
		return
	}

	switch d.Options.OutputFormat {
	case "csv":
		t.RenderCSV()
	case "html":
		t.RenderHTML()
	case "json":
		jsBytes, err := json.MarshalIndent(t, "", "  ")
		if err != nil {
			fmt.Printf("can't jsonify, err: %v", err)
			return
		}
		fmt.Printf("%s", jsBytes)
	case "markdown":
		t.RenderMarkdown()
	default:
		if !d.Options.DiffCalling {
			t.SetAutoIndex(true)
			t.Render()
		}
	}
}

func (db *kdbx) List() error {
	db.SelectedEntries = db.Entries
	db.Display()
	if db.Options.Quite {
		return nil
	}

	if !db.Options.NoKey {
		fmt.Printf("%sThis command used: \n\tkeyfile: %s\n\tdatabase: %s\n",
			utils.ColorGreen, db.Options.Key, db.Options.Database)
	} else {
		fmt.Printf("%sThis command used: \n\tdatabase: %s\n",
			utils.ColorGreen, db.Options.Database)
	}
	fmt.Printf("\nShowing %v of %v total entries%s\n",
		len(db.SelectedEntries), len(db.Entries), utils.ColorReset)
	return nil
}

func cacheFile(t table.Writer, cacheFilename string) {
	_, w, _ := os.Pipe()

	t.SetAutoIndex(false)
	t.SortBy([]table.SortBy{
		{Number: 1, Mode: table.Asc},
	})
	t.SetOutputMirror(w)
	data := t.RenderCSV() + "\n"
	w.Close()

	err := os.WriteFile(cacheFilename, []byte(data), 0600)
	if err != nil {
		log.Errorf("failed to write cache file: %v", err)
	}

	// log.Printf("wrote cachefile: %v for options: %#v", d.Options.CacheFile, d.Options.String())
}

func (d *kdbx) updateTableWithSelectedEntries(t table.Writer) table.Writer {
	if d.Options.Days > 0 {
		var newSelItems []models.Interested
		for _, ent := range d.SelectedEntries {
			if ent.Created.After(time.Now().AddDate(0, 0, -1*d.Options.Days)) ||
				ent.Modified.After(time.Now().AddDate(0, 0, -1*d.Options.Days)) {
				newSelItems = append(newSelItems, ent)
			}
		}
		d.SelectedEntries = newSelItems
	}

	for _, c := range d.SelectedEntries {
		cols := []string{c.Title, strconv.Itoa(c.Histories),
			c.Created.Format(TimeLayout), c.Modified.Format(TimeLayout)}
		switch d.Options.Fields {
		case "all":
			dataBytes, err := json.Marshal(c.KeyValues)
			if err != nil {
				fmt.Printf("marshall failed, err: %v", err)
			}
			cols = append(cols, string(dataBytes))
			t.AppendRow(table.Row{cols[0], cols[1], cols[2], cols[3], cols[4]})
		case "few":
			t.AppendRow(table.Row{cols[0]})
		default:
			t.AppendRow(table.Row{cols[0], cols[1], cols[2], cols[3]})
		}
	}
	t.SetStyle(table.StyleLight)
	return t
}

func (d *kdbx) getTable() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	switch d.Options.Fields {
	// header row
	case "all":
		t.AppendHeader(table.Row{"Title (col #1)", "Histories", "Created", "Modified", "Notes"})
	case "few":
		t.AppendHeader(table.Row{"Title (col #1)"})
	default:
		t.AppendHeader(table.Row{"Title (col #1)", "Histories", "Created", "Modified"})
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:     "Notes",
			WidthMax: 64,
			WidthMaxEnforcer: func(col string, maxLen int) string {
				if len(col) > 64 {
					return col[:61] + "..."
				}
				return col
			},
		},
	})
	return t
}
