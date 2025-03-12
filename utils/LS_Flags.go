package utils

import (
	"os"
	"strings"
)

// Args parses command-line arguments and extracts the directory path and flags.
//
// It scans through `os.Args`, identifying flags and setting corresponding boolean values:
//   - `-a` : Includes hidden files.
//   - `-R` : Enables recursive listing.
//   - `-t` : Sorts files by modification time.
//   - `-l` : Enables long listing format with detailed file information.
//   - `-r` : Reverses the sorting order.
//
// If an argument is not a flag (i.e., does not start with `-`), it is treated as the directory path.
//
// Returns:
//   - `path` (string): The specified directory path, or an empty string if none is provided (defaults to `.`).
//   - `lFlag`, `RFlag`, `aFlag`, `rFlag`, `tFlag` (bool): Boolean values representing whether each flag is set.
func Args() (path string, lFlag, RFlag, aFlag, rFlag, tFlag bool) {
	endOfFlags := false

	for _, arg := range os.Args[1:] {
		// If we encounter the double-dash, stop processing flags.
		if arg == "--" {
			endOfFlags = true
			continue
		}

		if !endOfFlags && strings.HasPrefix(arg, "-") && arg != "-" {
			if strings.Contains(arg, "a") {
				aFlag = true
			}
			if strings.Contains(arg, "R") {
				RFlag = true
			}
			if strings.Contains(arg, "t") {
				tFlag = true
			}
			if strings.Contains(arg, "l") {
				lFlag = true
			}
			if strings.Contains(arg, "r") {
				rFlag = true
			}
		} else {
			path = arg
		}
	}

	return
}
