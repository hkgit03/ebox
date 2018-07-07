package main

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

// sortableStringSlice makes []string sortable. The sorting is case-insensitive.
type sortableStringSlice []string

// implementation of sort.Interface on stringSlice
func (s sortableStringSlice) Len() int {
	return len(s)
}

// implementation of sort.Interface on stringSlice
func (s sortableStringSlice) Less(i, j int) bool {
	return strings.ToLower(s[i]) < strings.ToLower(s[j])
}

// implementation of sort.Interface on stringSlice
func (s sortableStringSlice) Swap(i, j int) {
	tmp := s[i]
	s[i] = s[j]
	s[j] = tmp
}

// checkNoEmacsBlocking checks that no running Emacs instance is using the
// current distribution. If such an instance is found, checkNoEmacsBlocking
// returns its pid and an error, else emptystring and nil.
func checkNoEmacsBlocking() (string, error) {

	// open /proc
	procDir, err := os.Open("/proc")
	if err != nil {
		panic(err)
	}

	// read all filenames from /proc
	filenames, err := procDir.Readdirnames(0)
	if err != nil {
		panic(err)
	}

	// read every /proc/*/cmdline file
	for _, pid := range filenames {

		// expected process name and command line switches to exclude a process
		// command line switches are separated with \x00 in the cmdline file
		emacsProcessName := "emacs"
		excludedSwitches := []string{"\x00-q", "\x00--no-init-file", "\x00-Q", "\x00--quick"}

		// read cmdline file
		cmdline, err := ioutil.ReadFile("/proc/" + pid + "/cmdline")
		if err != nil {
			// pid is probably not the directory of a pid
			continue
		}
		cmdString := string(cmdline)

		// check presence of excluded switches
		excludedSwitchPresent := false
		for _, excludedSwitch := range excludedSwitches {
			if strings.Contains(cmdString, excludedSwitch) {
				excludedSwitchPresent = true
				break
			}
		}

		// check process name and command line switches
		if strings.HasPrefix(cmdString, emacsProcessName) && !excludedSwitchPresent {
			return pid, errors.New("Blocking Emacs instance found")
		}
	}

	return "", nil
}
