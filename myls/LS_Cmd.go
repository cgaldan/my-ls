package myls

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	blue    = "\033[1;34m" // Bright Blue for directories
	green   = "\033[1;32m" // Bright Green for executable files
	cyan    = "\033[1;36m" // Bright Cyan for symbolic links
	red     = "\033[1;31m" // Bright Red for broken symbolic links
	magenta = "\033[1;35m" // Bright Magenta for image files
	yellow  = "\033[1;33m" // Bright Yellow
	black   = "\033[1;30m" // Bright Black
	reset   = "\033[0m"    // Resets to the default terminal color

	// Background colors
	bgBlack  = "\033[40m"
	bgRed    = "\033[41m"
	bgYellow = "\033[43m"
	bgBlue   = "\033[44m"
)

type MyLSFiles struct {
	Name           string
	IsDir          bool
	IsExec         bool
	IsLink         bool
	IsBroken       bool
	IsBlockDevice  bool
	IsCharDevice   bool
	IsSocket       bool
	IsPipe         bool
	IsOrphanedLink bool
	IsSetuid       bool
	IsSetgid       bool
	IsStickyDir    bool
	Size           int64
	ModTime        time.Time
	Mode           fs.FileMode
	OwnerName      string
	GroupName      string
	NLink          uint64
}

// GetColor returns the appropriate ANSI color code based on the file type.
//   - Directories are displayed in blue.
//   - Executable files are displayed in green.
//   - Symbolic links are displayed in cyan.
//   - Broken symbolic links are displayed in red.
//   - Regular files are displayed in the default terminal color.
func (file MyLSFiles) GetColor() string {
	ext := strings.ToLower(filepath.Ext(file.Name))

	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".mp4":
		return magenta
	case ".mp3", ".wav", ".ogg", ".flac":
		return cyan
	case ".zip", ".tar", ".gz", ".bz2", ".rar", ".7z", ".deb", ".rpm":
		return red
	}
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
	if file.IsBlockDevice || file.IsCharDevice || file.IsPipe {
		return bgBlack + yellow
	}
	if file.IsSocket {
		return magenta
	}
	if file.IsSetuid {
		return bgRed
	}
	if file.IsSetgid {
		return bgYellow + black
	}
	if file.IsStickyDir {
		return bgBlue
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

	fileInfo, err := os.Stat(dirName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if !fileInfo.IsDir() {
		printFileDetails(dirName, fileInfo, lFlag)
		return
	}

	entries, err := os.ReadDir(dirName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var totalBlocks int64
	var maxSize, maxNlink int

	if aFlag {
		dotFile := getFileAttributes(".", fileInfo)
		files = append(files, dotFile)
		if stat, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
			totalBlocks += int64(stat.Blocks)
		}

		if parentInfo, err := os.Stat(".."); err == nil {
			parentFile := getFileAttributes("..", parentInfo)
			files = append(files, parentFile)
			if stat, ok := parentInfo.Sys().(*syscall.Stat_t); ok {
				totalBlocks += int64(stat.Blocks)
			}
			updateMaxLengths(&maxNlink, &maxSize, parentFile)
		}
	}

	for _, entry := range entries {
		fileName := entry.Name()
		if !aFlag && strings.HasPrefix(fileName, ".") {
			continue
		}

		info, err := os.Lstat(filepath.Join(dirName, fileName))
		if err != nil {
			continue
		}

		file := getFileAttributes(filepath.Join(dirName, fileName), info)
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			totalBlocks += int64(stat.Blocks)
		}
		files = append(files, file)
		if RFlag && file.IsDir {
			subDirs = append(subDirs, filepath.Join(dirName, file.Name))
		}
		updateMaxLengths(&maxNlink, &maxSize, file)
	}

	sortFiles(&files, tFlag, rFlag)

	if RFlag {
		fmt.Println(dirName + ":")
	}

	if lFlag {
		fmt.Printf("total %d\n", totalBlocks/2)
		for _, file := range files {
			fmt.Println(formatLongEntry(file, maxNlink, maxSize))
		}
	} else {
		printFiles(files)
	}

	if RFlag {
		for _, subDir := range subDirs {
			fmt.Println()
			TheMainLS(subDir, lFlag, RFlag, aFlag, rFlag, tFlag)
		}
	}
}

func printFileDetails(path string, info os.FileInfo, lFlag bool) {
	file := getFileAttributes(path, info)
	if lFlag {
		fmt.Println(formatLongEntry(file, len(strconv.Itoa(int(file.NLink))), len(strconv.Itoa(int(file.Size)))))
	} else {
		printFiles([]MyLSFiles{file})
	}
}

func updateMaxLengths(maxNlink, maxSize *int, file MyLSFiles) {
	strNLink := strconv.Itoa(int(file.NLink))
	strSize := strconv.Itoa(int(file.Size))
	if *maxNlink < len(strNLink) {
		*maxNlink = len(strNLink)
	}
	if *maxSize < len(strSize) {
		*maxSize = len(strSize)
	}
}

func sortFiles(files *[]MyLSFiles, tFlag, rFlag bool) {
	if tFlag {
		sortByTime(*files)
	} else {
		sortByName(*files)
	}
	if rFlag {
		reverseFiles(*files)
	}
}

func getFileAttributes(path string, info os.FileInfo) MyLSFiles {
	stat, _ := info.Sys().(*syscall.Stat_t)
	var nlink uint64 = 1
	var uid, gid uint32

	if stat != nil {
		nlink = uint64(stat.Nlink)
		uid = stat.Uid
		gid = stat.Gid
	}

	ownerName := strconv.Itoa(int(uid))
	if owner, err := user.LookupId(strconv.Itoa(int(uid))); err == nil {
		ownerName = owner.Username
	}

	groupName := strconv.Itoa(int(gid))
	if group, err := user.LookupGroupId(strconv.Itoa(int(gid))); err == nil {
		groupName = group.Name
	}

	return MyLSFiles{
		Name:          filepath.Base(path),
		IsDir:         info.IsDir(),
		IsExec:        !info.IsDir() && (info.Mode().Perm()&0o111 != 0),
		IsLink:        info.Mode()&os.ModeSymlink != 0,
		IsBroken:      info.Mode()&os.ModeSymlink != 0 && !exists(path),
		IsBlockDevice: info.Mode()&os.ModeDevice != 0 && info.Mode()&syscall.S_IFBLK != 0,
		IsCharDevice:  info.Mode()&os.ModeDevice != 0 && info.Mode()&syscall.S_IFCHR != 0,
		IsSocket:      info.Mode()&os.ModeSocket != 0,
		IsPipe:        info.Mode()&os.ModeNamedPipe != 0,
		IsSetuid:      info.Mode()&os.ModeSetuid != 0,
		IsSetgid:      info.Mode()&os.ModeSetgid != 0,
		IsStickyDir:   info.IsDir() && info.Mode()&os.ModeSticky != 0,
		Size:          info.Size(),
		ModTime:       info.ModTime(),
		Mode:          info.Mode(),
		OwnerName:     ownerName,
		GroupName:     groupName,
		NLink:         nlink,
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
