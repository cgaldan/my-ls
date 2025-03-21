package utils

import (
	"fmt"
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
func Args() (paths []string, lFlag, RFlag, aFlag, rFlag, tFlag bool) {
	endOfFlags := false
	flags := "-aRtlr"

	for _, arg := range os.Args[1:] {

		if arg == "--help" {
			fmt.Println("Usage: ./myls [OPTION]... [FILE]...")
			fmt.Println("List information about the FILEs (the current directory by default).")
			fmt.Println("Options:")
			fmt.Println("  -a  : Includes hidden files.")
			fmt.Println("  -R  : Enables recursive listing.")
			fmt.Println("  -t  : Sorts files by modification time.")
			fmt.Println("  -l  : Enables long listing format with detailed file information.")
			fmt.Println("  -r  : Reverses the sorting order.")
			os.Exit(0)
		}

		// If we encounter the double-dash, stop processing flags.
		if arg == "--" {
			endOfFlags = true
			continue
		}

		if !endOfFlags && strings.HasPrefix(arg, "-") && arg != "-" {
			for _, r := range arg {
				if !strings.ContainsAny(flags, string(r)) {
					fmt.Printf("myls: invalid option -- '%s'\n", string(r))
					fmt.Println("Try '--help' flag for more information.")
					os.Exit(0)
				}
			}
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
			paths = append(paths, arg)
		}
	}

	return
}
