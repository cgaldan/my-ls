package sortpkg

import (
	"ls/data"
	"os"
	"strings"
)

func SortFiles(files *[]data.MyLSFiles, tFlag, rFlag bool) {
	var caseSensitive bool

	if tFlag {
		sortByTime(*files, caseSensitive)
	} else {
		caseSensitive = isCaseSensitiveSort()
		if caseSensitive {
			sortByNameCaseSensitive(*files)
		} else {
			sortByName(*files)
		}

	}

	if rFlag {
		reverseFiles(*files)
	}
}

// normalizeASCII keeps printable ASCII characters and converts A-Z to a-z.
func normalizeASCII(name string) string {
	var b strings.Builder

	isPunctuation := true

	for i := 0; i < len(name); i++ {
		c := name[i]
		if !(c >= ' ' && c <= '~') || (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			isPunctuation = false
			break
		}
	}
	if isPunctuation {
		return name
	}

	for i := 0; i < len(name); i++ {
		c := name[i]
		switch {
		case c >= '0' && c <= '9':
			b.WriteByte(c)
		case c >= 'a' && c <= 'z':
			b.WriteByte(c)
		case c >= 'A' && c <= 'Z':
			// convert uppercase to lowercase
			b.WriteByte(c + ('a' - 'A'))
			// case c >= ' ' && c <= '~': // keep printable ASCII characters including punctuation
			// 	b.WriteByte(c)
		}
	}
	return b.String()
}

func sortByName(files []data.MyLSFiles) {
	n := len(files)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if normalizeASCII(files[j].Name) > normalizeASCII(files[j+1].Name) {
				files[j], files[j+1] = files[j+1], files[j]
			}
		}
	}
}

func sortByNameCaseSensitive(files []data.MyLSFiles) {
	n := len(files)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if files[j].Name > files[j+1].Name {
				files[j], files[j+1] = files[j+1], files[j]
			}
		}
	}
}

// sortByTime sorts the given slice of MyLSFiles by modification time in descending order,
// so that the most recently modified files appear first.
// It also uses a bubble sort algorithm for ordering.
func sortByTime(files []data.MyLSFiles, caseSensitive bool) {
	n := len(files)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if files[j].ModTime.Before(files[j+1].ModTime) {
				files[j], files[j+1] = files[j+1], files[j]
			} else if files[j].ModTime.Equal(files[j+1].ModTime) {
				a, b := files[j].Name, files[j+1].Name
				if caseSensitive {
					a = normalizeASCII(a)
					b = normalizeASCII(b)
				}
				if a > b {
					files[j], files[j+1] = files[j+1], files[j]
				}
			}
		}
	}
}

// reverseFiles reverses the order of the given slice of MyLSFiles.
// This is useful when the `-r` flag is enabled to display results in reverse order.
func reverseFiles(files []data.MyLSFiles) {
	n := len(files)
	for i := 0; i < n/2; i++ {
		files[i], files[n-1-i] = files[n-1-i], files[i]
	}
}

func isCaseSensitiveSort() bool {
	// 1) pick up LC_COLLATE (preferred) or LANG
	loc := os.Getenv("LC_COLLATE")
	if loc == "" {
		loc = os.Getenv("LANG")
	}

	// 2) lowercase & strip off any ".UTF-8" or "@modifier"
	loc = strings.ToLower(loc)
	if i := strings.IndexAny(loc, ".@"); i != -1 {
		loc = loc[:i]
	}

	// 3) only the special C/POSIX locales are case-sensitive
	return loc == "c" || loc == "posix"
}
