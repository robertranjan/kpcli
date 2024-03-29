package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
)

// NewDiff returns a *Diff
func NewDiff(opts Options) *Diff {
	var fromDBOpts = Options{
		CacheFile:    "database2.out",
		Database:     opts.Database2,
		Days:         10000,
		DiffCalling:  true,
		Fields:       "few",
		Key:          opts.Key,
		Notify:       true,
		OutputFormat: opts.OutputFormat,
		Pass:         opts.Pass,
		Reverse:      true,
		// Quite:        true,
		// options for diff cmd
	}
	var toDBOpts = Options{
		CacheFile:    "database1.out",
		Database:     opts.Database,
		Days:         10000,
		DiffCalling:  true,
		Fields:       "few",
		Key:          opts.Key,
		Notify:       true,
		OutputFormat: opts.OutputFormat,
		Pass:         opts.Pass,
		Reverse:      true,
		// Quite:        true,
		// options for diff cmd
	}

	return &Diff{
		ToDBOption:   &toDBOpts,
		FromDBOption: &fromDBOpts,
		options:      &opts,
	}
}

func getRecentFile(dir string, filename string) string {
	pattern := fmt.Sprintf("%v.*", filename)
	var recent time.Time
	var recentFile string

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, file := range files {
		fi, err := file.Info()
		if err != nil {
			return ""
		}
		matched, err := regexp.MatchString(pattern, file.Name())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if matched && fi.ModTime().After(recent) {
			recent = fi.ModTime()
			recentFile = file.Name()
		}
	}
	fmt.Printf("recentfile: %s\n", filepath.Join(dir, recentFile))
	return filepath.Join(dir, recentFile)
}

// Diff shows the difference between 2 databases
// notify option can be used to notify your email id (work only for gmail at the moment)
func (d *Diff) Diff() error {

	// list entries from recent backup
	dbOne, err := NewDB(*d.FromDBOption)
	if err != nil {
		return err
	}
	err = dbOne.List()
	if err != nil {
		return err
	}

	// list current db entries
	dbTwo, err := NewDB(*d.ToDBOption)
	if err != nil {
		return err
	}

	err = dbTwo.List()
	if err != nil {
		return err
	}

	outputHeader := []byte(fmt.Sprintf("here are the diffs between %v and %v\n",
		d.options.Database2, d.options.Database))
	cmd := exec.Command("diff", []string{"database2.out", "database1.out"}...)
	outputHeader = append(outputHeader, []byte(strings.Repeat("-", 70))...)
	outputHeader = append(outputHeader, []byte("\n")...)
	out, _ := cmd.CombinedOutput()
	if string(out) == "" {
		out = []byte("\tNo differences")
	}
	outputLines := strings.Split(string(out), "\n")
	sort.Strings(outputLines)

	var ANSILines, HTMLLines []string
	var htmlLine, ansiLine string
	for _, line := range outputLines {
		if strings.HasPrefix(line, "<") || strings.HasPrefix(line, ">") {

			// color code output for terminal
			ansiLine = strings.Replace(line, "<", "\033[31m (removed) \033[0m", -1)
			ansiLine = strings.Replace(ansiLine, ">", "\033[32m ( added ) \033[0m", -1)
			ANSILines = append(ANSILines, ansiLine)

			// color code output for html email 'from above var: ansiLine'
			//   to avoid repl issues with html tags[< & >]
			htmlLine = strings.Replace(ansiLine, "\033[31m (removed) \033[0m",
				"<font color=red> (removed) </font>", -1)
			htmlLine = strings.Replace(htmlLine, "\033[32m ( added ) \033[0m",
				"<font color=green> (added) </font>", -1)
			HTMLLines = append(HTMLLines, htmlLine)
		}
	}

	fmt.Printf("%s", outputHeader)
	fmt.Printf("%v", strings.Join(ANSILines, "\n"))

	HTMLOut := []byte("<pre>")
	HTMLOut = append(HTMLOut, outputHeader...)
	HTMLOut = append(HTMLOut, []byte(strings.Join(HTMLLines, "\n"))...)
	HTMLOut = append(HTMLOut, []byte("</pre>")...)
	err = os.WriteFile(d.options.OutputFilename, HTMLOut, 0600)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	if d.options.Notify && len(ANSILines) > 0 {
		d.Notify(d.options.OutputFilename)
	} else {
		color.Yellow("\n  >>> Not sending any emails " +
			"as there is no changes or notification wasn't requested.\n")
	}

	// fmt.Printf("fromDB: %v\ntoDB: %v\nopts: %v\n", d.FromDBOption, d.ToDBOption, d.options)

	return nil
}
