package myls

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

const (
	blue  = "\033[1;34m" // Bright Blue for directories
	green = "\033[1;32m" // Bright Green for executable files
	cyan  = "\033[1;36m" // Bright Cyan for symbolic links
	red   = "\033[1;31m" // Bright Red for broken symbolic links
	reset = "\033[0m"    // Resets to the default terminal color
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

// GetColor returns the appropriate ANSI color code based on the file type.
//   - Directories are displayed in blue.
//   - Executable files are displayed in green.
//   - Symbolic links are displayed in cyan.
//   - Broken symbolic links are displayed in red.
//   - Regular files are displayed in the default terminal color.
func (file MyLSFiles) GetColor() string {
	if file.IsBroken {
		return red
	}
	if file.IsLink {
		return cyan
	}
	if file.IsDir {
		return blue
	}
	if file.IsExec {
		return green
	}
	return reset
}

// TheMainLS lists directory contents similar to the Unix `ls` command.
// It supports various flags for additional functionality:
//
//   - `lFlag` : Enables long listing format with detailed file information.
//   - `RFlag` : Recursively lists subdirectories.
//   - `aFlag` : Includes hidden files (those starting with `.`).
//   - `rFlag` : Reverses the sorting order.
//   - `tFlag` : Sorts files by modification time (newest first).
//
// Parameters:
//   - `dirName` (string): The directory to list. Defaults to the current directory if empty.
//   - `lFlag`, `RFlag`, `aFlag`, `rFlag`, `tFlag` (bool): Flags controlling the behavior.
//
// The function retrieves directory contents, filters them based on flags, sorts them,
// and prints the results with color coding. If `RFlag` is set, it recursively lists subdirectories.
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

	var totalBlocks int64 = 0

	for _, entry := range entries {
		fileName := entry.Name()
		if !aFlag && strings.HasPrefix(fileName, ".") {
			continue
		}
		info, err := os.Lstat(dirName + "/" + fileName)
		if err != nil {
			continue
		}
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			totalBlocks += int64(stat.Blocks)
		}

		isExec := !info.IsDir() && (info.Mode().Perm()&0o111 != 0)
		isLink := info.Mode()&os.ModeSymlink != 0
		isBroken := isLink && !exists(dirName+"/"+fileName)

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

	if lFlag {
		fmt.Printf("total %d\n", totalBlocks/2)
		for _, file := range files {
			fullPath := filepath.Join(dirName, file.Name)
			info, err := os.Lstat(fullPath)
			if err != nil {
				continue
			}
			fmt.Println(formatLongEntry(file, info))
		}
		return
	}

	if RFlag {
		for _, file := range files {
			if file.IsDir {
				subDirs = append(subDirs, dirName+"/"+file.Name)
			}
		}
		fmt.Println(dirName + ":")
	}

	printFiles(files)
	// for _, file := range files {
	// 	fileName := formatFileNames(file.Name)
	// 	fmt.Print(file.GetColor() + fileName + reset + "  ")
	// }
	// fmt.Println()

	if RFlag {
		for _, subDir := range subDirs {
			fmt.Println()
			TheMainLS(subDir, lFlag, RFlag, aFlag, rFlag, tFlag)
		}
	}
}

// sortByName sorts the given slice of MyLSFiles in ascending order by name (case-insensitive).
// It uses a simple bubble sort algorithm for ordering.
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

// sortByTime sorts the given slice of MyLSFiles by modification time in descending order,
// so that the most recently modified files appear first.
// It also uses a bubble sort algorithm for ordering.
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

// reverseFiles reverses the order of the given slice of MyLSFiles.
// This is useful when the `-r` flag is enabled to display results in reverse order.
func reverseFiles(files []MyLSFiles) {
	n := len(files)
	for i := 0; i < n/2; i++ {
		files[i], files[n-1-i] = files[n-1-i], files[i]
	}
}

// exists checks whether a file or directory exists at the given path.
// Returns true if the file exists, otherwise returns false.
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
