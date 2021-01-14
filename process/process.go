package process

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	sPrefix = "php-fpm: pool"
)

var (
	rePoolName = regexp.MustCompile(`^php-fpm: pool (.*)$`)
)

// Gathers all processes in /proc and tally all PHP-FPM pool workers
func TallyPHPPools() (map[string]uint64, error) {
	pools := make(map[string]uint64)

	// Find all procs
	files, err := ioutil.ReadDir(procPath())
	if err != nil {
		return pools, err
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		cmd, err := procCmd(file.Name())
		if err != nil {
			continue
		}

		if !strings.HasPrefix(cmd, sPrefix) {
			continue
		}

		// Get pool part
		pool := rePoolName.ReplaceAllString(strings.TrimRight(cmd, "\x00"), "$1")

		if _, ok := pools[pool]; ok {
			pools[pool]++
		} else {
			pools[pool] = 1
		}
	}

	return pools, nil
}

// Build the path for a file in /proc
func procPath(parts ...string) string {
	parts = append([]string{"/proc"}, parts...)
	return filepath.Join(parts...)
}

func procCmd(pid string) (string, error) {
	cmdPath := procPath(pid, "cmdline")
	contents, err := ioutil.ReadFile(cmdPath)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}
