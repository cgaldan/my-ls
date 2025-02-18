package myls

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"
)

const (
	blue  = "\033[1;34m" // Bright Blue για φακέλους
	green = "\033[1;32m" // Bright Green για εκτελέσιμα
	cyan  = "\033[1;36m" // Bright Cyan για symlinks
	red   = "\033[1;31m" // Bright Red για broken symlinks
	reset = "\033[0m"    // Επαναφορά στο default
)

type MyLSFiles struct {
	Name     string
	IsDir    bool
	IsExec   bool
	IsLink   bool
	IsBroken bool
	Size     int64
	ModTime  time.Time
	Mode     fs.FileMode
}

func (file MyLSFiles) GetColor() string {
	if file.IsBroken {
		return red // Σπασμένο symlink
	}
	if file.IsLink {
		return cyan // Συμβολικός σύνδεσμος
	}
	if file.IsDir {
		return blue // Φάκελος
	}
	if file.IsExec {
		return green // Εκτελέσιμο αρχείο
	}
	return reset // Κανονικό αρχείο
}

func TheMainLS(dirName string, lFlag, RFlag, aFlag, rFlag, tFlag bool) {

	var files []MyLSFiles
	var subDirs []string

	if dirName == "" {
		dirName = "."
	}

	var myFS fs.FS = os.DirFS(dirName)
	entries, err := fs.ReadDir(myFS, ".")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if aFlag {

		if dotInfo, err := os.Lstat(dirName); err == nil {
			files = append(files, MyLSFiles{
				Name:    ".",
				IsDir:   dotInfo.IsDir(),
				IsExec:  (dotInfo.Mode().Perm() & 0o111) != 0,
				IsLink:  (dotInfo.Mode() & os.ModeSymlink) != 0,
				Size:    dotInfo.Size(),
				ModTime: dotInfo.ModTime(),
				Mode:    dotInfo.Mode(),
			})
		}

		if parentInfo, err := os.Lstat("../"); err == nil {
			files = append(files, MyLSFiles{
				Name:    "..",
				IsDir:   parentInfo.IsDir(),
				IsExec:  (parentInfo.Mode().Perm() & 0o111) != 0,
				IsLink:  (parentInfo.Mode() & os.ModeSymlink) != 0,
				Size:    parentInfo.Size(),
				ModTime: parentInfo.ModTime(),
				Mode:    parentInfo.Mode(),
			})
		}
	}

	for _, entry := range entries {
		fileName := entry.Name()
		if !aFlag && strings.HasPrefix(fileName, ".") {
			continue
		}
		info, err := os.Lstat(dirName + "/" + fileName) // Χρήση `Lstat` για symlinks
		if err != nil {
			continue
		}

		isExec := !info.IsDir() && (info.Mode().Perm()&0o111 != 0)
		isLink := info.Mode()&os.ModeSymlink != 0
		isBroken := isLink && !exists(dirName+"/"+fileName) // Broken symlink

		files = append(files, MyLSFiles{
			Name:     fileName,
			IsDir:    entry.IsDir(),
			IsExec:   isExec,
			IsLink:   isLink,
			IsBroken: isBroken,
			Size:     info.Size(),
			ModTime:  info.ModTime(),
			Mode:     info.Mode(),
		})

	}

	if tFlag {
		sortByTime(files)
	} else {
		sortByName(files)
	}

	if rFlag {
		reverseFiles(files)
	}

	if RFlag {
		for _, file := range files {
			if file.IsDir {
				subDirs = append(subDirs, dirName+"/"+file.Name)
			}
		}
		fmt.Println("\n" + dirName + ":")
	}
	for _, file := range files {
		fmt.Print(file.GetColor() + file.Name + reset + "  ")
	}
	fmt.Println()

	if RFlag {
		for _, subDir := range subDirs {
			TheMainLS(subDir, lFlag, RFlag, aFlag, rFlag, tFlag)
		}
	}
}

func sortByName(files []MyLSFiles) {
	n := len(files)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if strings.ToLower(files[j].Name) > strings.ToLower(files[j+1].Name) {
				files[j], files[j+1] = files[j+1], files[j]
			}
		}
	}
}

func sortByTime(files []MyLSFiles) {
	n := len(files)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if files[j].ModTime.Before(files[j+1].ModTime) {
				files[j], files[j+1] = files[j+1], files[j]
			}
		}
	}
}

func reverseFiles(files []MyLSFiles) {
	n := len(files)
	for i := 0; i < n/2; i++ {
		files[i], files[n-1-i] = files[n-1-i], files[i]
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
