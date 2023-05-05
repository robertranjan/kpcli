package ls

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/tobischo/gokeepasslib/v3"
)

var (
	colorGreen = "\033[32m"
	colorReset = "\033[0m"

	TimeLayout = "2006-01-02 15:04:05"
	lengthUser = 25
)

// NewDB create and return a new kdbx db object
func NewDB(opts Options) (*db, error) {

	return &db{
		Options: &opts,
	}, nil
}

func (d *db) AppendEntries(grp gokeepasslib.Group) {
	for _, e := range grp.Entries {
		kv := make(map[string]string)
		for _, entry := range e.Values {
			kv[entry.Key] = entry.Value.Content
		}

		var hist int
		if len(e.Histories) > 0 {
			hist = len(e.Histories[0].Entries)
		}
		et := Interested{
			Title:     grp.Name + "/" + strings.TrimSpace(e.GetTitle()),
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
}

func (d *db) loadEntries() {
	file, err := os.Open(d.Options.Database)
	if err != nil {
		log.Fatalf("failed open database file: %v, err: %v", d.Options.Database, err)
	}

	db := gokeepasslib.NewDatabase(gokeepasslib.WithDatabaseKDBXVersion4())
	cred, err := gokeepasslib.NewPasswordAndKeyCredentials(d.Options.Pass, d.Options.Key)
	if err != nil {
		log.Fatalf("failed to create credentials, err: %v", err)
	}
	db.Credentials = cred

	if err := gokeepasslib.NewDecoder(file).Decode(db); err != nil {
		log.Fatalf("failed to decode db: %#v, err: %v", db, err)
	}

	if err := db.UnlockProtectedEntries(); err != nil {
		log.Fatalf("failed to unload db, err: %v", err)
	}

	for _, rootgp := range db.Content.Root.Groups {
		for _, grp := range rootgp.Groups {
			d.AppendEntries(grp)
		}
		d.AppendEntries(rootgp)
	}
}

func (d *db) updateSelectedEntries() {
	if d.Options.Days == 0 {
		return
	}
	var newSelItems []Interested
	for _, ent := range d.SelectedEntries {
		if ent.Created.After(time.Now().AddDate(0, 0, -1*d.Options.Days)) ||
			ent.Modified.After(time.Now().AddDate(0, 0, -1*d.Options.Days)) {
			newSelItems = append(newSelItems, ent)
		}
	}
	d.SelectedEntries = newSelItems
}

func (d *db) getTableHeader() table.Writer {
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
	return t
}

func (d *db) display() {

	t := d.getTableHeader()
	d.updateSelectedEntries()
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
	// t.SetPageSize(5)

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

	// Note:
	//    it sort with string sort. Try `--sort-by-col 2`` to see how it sorts.
	if d.Options.Reverse {
		t.SortBy([]table.SortBy{
			{Number: d.Options.SortbyCol, Mode: table.Dsc},
		})
	} else {
		t.SortBy([]table.SortBy{
			{Number: d.Options.SortbyCol, Mode: table.Asc},
		})
	}

	fmt.Println()
	_, w, _ := os.Pipe()

	// when cacheFile is set
	//   `./kpcli diff` uses this to save output in file & then diff them
	if d.Options.CacheFile != "" {
		t.SetAutoIndex(false)
		t.SortBy([]table.SortBy{
			{Number: 1, Mode: table.Asc},
		})
		t.SetOutputMirror(w)
		data := t.RenderCSV() + "\n"
		w.Close()

		os.WriteFile(d.Options.CacheFile, []byte(data), 0600)
		// log.Printf("wrote cachefile: %v for options: %#v", d.Options.CacheFile, d.Options.String())
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
		if !d.Options.Diff {
			t.SetAutoIndex(true)
			t.Render()
		}
	}
}

func (d *db) show() {

	d.SelectedEntries = d.Entries
	d.display()
	if d.Options.Quite {
		return
	}
	fmt.Printf("\n\n%sThis command used: \n\tkeyfile: %s\n\tdatabase: %s\n",
		colorGreen, d.Options.Key, d.Options.Database)
	fmt.Printf("\nShowing %v of %v total entries%s\n",
		len(d.SelectedEntries), len(d.Entries), colorReset)
}

// List show the entries from a kdbx database
// You may call it to save the listing to file by specifying right options
func (d *db) List() error {
	d.loadEntries()
	d.show()
	// fmt.Printf("option: %#v\n", d.Options)
	return nil
}
