package diff

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/robertranjan/kpcli/cmds/ls"
)

func NewDiff(opts Options) *Diff {
	var fromDBOpts = ls.Options{
		Reverse:      true,
		Days:         10000,
		Pass:         opts.Pass,
		Database:     opts.Database2,
		Key:          opts.Key,
		Fields:       "few",
		CacheFile:    "database2.out",
		Quite:        true,
		Diff:         true,
		OutputFormat: opts.OutputFormat,
	}
	var toDBOpts = ls.Options{
		Reverse:      true,
		Days:         10000,
		Pass:         opts.Pass,
		Database:     opts.Database,
		Key:          opts.Key,
		Fields:       "few",
		CacheFile:    "database1.out",
		Quite:        true,
		Diff:         true,
		OutputFormat: opts.OutputFormat,
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
	return filepath.Join(dir, recentFile)
}

func (d *Diff) Run() error {
	// list current db entries
	lsCurDB, err := ls.NewDB(*d.ToDBOption)
	if err != nil {
		return err
	}
	var errMsg string
	err = lsCurDB.Run()
	if err != nil {
		errMsg = err.Error()
	}

	// list recent backedup db entries
	lsOldDB, err := ls.NewDB(*d.FromDBOption)
	if err != nil {
		return err
	}
	err = lsOldDB.Run()
	if err != nil {
		errMsg += err.Error()
	}

	outputHeader := []byte(fmt.Sprintf("Running diff between \n\t%v and \n\t%v\n", d.options.Database2, d.options.Database))
	outputHeader = append(outputHeader, []byte("\nhere are the diffs:\n")...)
	cmd := exec.Command("diff", []string{"database2.out", "database1.out"}...)
	moreHeader := []byte(fmt.Sprintf(" %s to %v\n %s\n",
		filepath.Base(d.options.Database2), filepath.Base(d.options.Database),
		strings.Repeat("-", 70)),
	)
	outputHeader = append(outputHeader, moreHeader...)
	out, _ := cmd.CombinedOutput()
	if string(out) == "" {
		out = []byte("\tNo differneces")
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

			// color code output for html email 'from above var: ansiLine' to avoid repl issues with html tags[< & >]
			htmlLine = strings.Replace(ansiLine, "\033[31m (removed) \033[0m", "<font color=red> (removed) </font>", -1)
			htmlLine = strings.Replace(htmlLine, "\033[32m ( added ) \033[0m", "<font color=green> (added) </font>", -1)
			HTMLLines = append(HTMLLines, htmlLine)
		}
	}

	fmt.Printf("%s", outputHeader)
	fmt.Printf("%v", strings.Join(ANSILines, "\n"))

	HTMLOut := []byte("<pre>")
	HTMLOut = append(HTMLOut, outputHeader...)
	HTMLOut = append(HTMLOut, []byte(strings.Join(HTMLLines, "\n"))...)
	HTMLOut = append(HTMLOut, []byte("</pre>")...)
	os.WriteFile(d.options.OutputFilename, HTMLOut, 0600)

	if d.options.Notify && len(ANSILines) > 0 {
		d.Notify(d.options.OutputFilename)
	} else {
		color.Yellow("\n  >>> Not sending any emails. Either no change or notifications wasn't requested.\n")
	}

	if errMsg != "" {
		return errors.New(errMsg)
	}
	return nil
}