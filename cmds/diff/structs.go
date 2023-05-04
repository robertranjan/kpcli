package diff

import "github.com/robertranjan/kpcli/cmds/ls"

type Options struct {
	Database       string
	Database2      string
	Key            string
	Notify         bool
	OutputFilename string
	OutputFormat   string
	Pass           string
}

type Diff struct {
	ToDBOption   *ls.Options
	FromDBOption *ls.Options
	options      *Options
}
