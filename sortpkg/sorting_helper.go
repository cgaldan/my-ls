package sortpkg

import (
	"ls/data"
	"os"
	"strings"
)

// var insensitiveLocales = map[string]bool{
// 	"en_us.utf8": true,
// 	"en_gb.utf8": true,
// 	"de_de.utf8": true,
// 	"fr_fr.utf8": true,
// 	"es_es.utf8": true,
// 	"it_it.utf8": true,
// 	"nl_nl.utf8": true,
// 	"pt_br.utf8": true,
// 	"pt_pt.utf8": true,
// 	"ru_ru.utf8": true,
// 	"zh_cn.utf8": true,
// 	"zh_tw.utf8": true,
// 	"ja_jp.utf8": true,
// }

func SortFiles(files *[]data.MyLSFiles, tFlag, rFlag bool) {
	var caseSensitive bool

	if tFlag {
		sortByTime(*files)
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

// sortByName sorts the given slice of MyLSFiles in ascending order by name (case-insensitive).
// It uses a simple bubble sort algorithm for ordering.
func sortByName(files []data.MyLSFiles) {
	n := len(files)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if strings.ToLower(files[j].Name) > strings.ToLower(files[j+1].Name) {
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
func sortByTime(files []data.MyLSFiles) {
	n := len(files)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if files[j].ModTime.Before(files[j+1].ModTime) {
				files[j], files[j+1] = files[j+1], files[j]
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

// func isCaseSensitiveSort() bool {
// 	locale := strings.ToLower(os.Getenv("LC_COLLATE"))
// 	if locale == "" {
// 		locale = strings.ToLower(os.Getenv("LANG"))
// 	}

// 	if insensitiveLocales[locale] {
// 		return false
// 	}

// 	if strings.HasPrefix(locale, "C") || strings.HasPrefix(locale, "POSIX") {
// 		return true
// 	}

// 	return true
// }

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
