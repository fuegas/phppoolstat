package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fuegas/phppoolstat/process"
	"github.com/fuegas/phppoolstat/utils"
)

const (
	sUsage = `phppoolstat, tally PHP-FPM pool workers for Telegraf

Usage:

  phppoolstat [options]

The options are:

  --tag <name>=<value>    Add a tag to the output (for example: env=production)

  --help                  Show this description
  --version               Show the current version
`
)

type ArrayFlags []string

var (
	// Git related variables
	branch  string
	commit  string
	version string

	// CMD flags
	fTags    ArrayFlags
	fVersion = flag.Bool("version", false, "Show the current version")
)

func main() {
	var err error

	// Parse arguments
	flag.Usage = func() { exitWithUsage(0) }
	flag.Var(&fTags, "tag", "tag to add to the output")
	flag.Parse()

	// Check that no unknown flags were passed
	args := flag.Args()
	if len(args) > 0 {
		utils.PrintError("Unknown options passed:", args)
		exitWithUsage(1)
	}

	// Show version if requested
	if *fVersion {
		if version == "" {
			fmt.Printf("experimental @ %s on %s\n", commit, branch)
		} else {
			fmt.Printf("v%s\n", version)
		}
		return
	}

	// Determine tags
	tagsMap := make(map[string]string)
	for _, tagFlag := range fTags {
		arr := strings.SplitN(tagFlag, "=", 2)
		if len(arr) != 2 {
			utils.PrintError("Provided tag contains more than one = character: ", tagFlag)
			exitWithUsage(1)
		}

		tagsMap[arr[0]] = arr[1]
	}

	tags := new(bytes.Buffer)
	for key, value := range tagsMap {
		fmt.Fprintf(tags, ",%s=%s", key, value)
	}

	// Storage for pool information
	pools, err := process.TallyPHPPools()
	if err != nil {
		exitError("Error trying to tally PHP-FPM pools:", err)
	}

	// Total number of pools
	total := uint64(0)

	// Build output prefix
	prefix := fmt.Sprintf("phpfpm_pools%s", tags.String())

	// Output pools
	for pool, count := range pools {
		fmt.Printf("%s,pool=%s count=%di\n", prefix, utils.Escape(pool), count)
		total = total + count
	}

	// Output total count of pools
	fmt.Printf("%s,pool=_all_ count=%di\n", prefix, total)
}

// Show the error message and exit
func exitError(msgs ...interface{}) {
	utils.PrintError(msgs...)
	os.Exit(1)
}

// Show usage message and exit with the provided code
func exitWithUsage(code int) {
	fmt.Println(sUsage)
	os.Exit(code)
}

// Set defaults if no values were passed
func init() {
	if commit == "" {
		commit = "unknown"
	}
	if branch == "" {
		branch = "unknown"
	}
}

func (p *ArrayFlags) String() string {
	return "string representation"
}

func (p *ArrayFlags) Set(value string) error {
	*p = append(*p, value)
	return nil
}
