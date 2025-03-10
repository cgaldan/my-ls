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
		info, err := os.Lstat(dirName)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		isExec := !info.IsDir() && (info.Mode().Perm()&0o111 != 0)
		isLink := info.Mode()&os.ModeSymlink != 0
		isBroken := isLink && !exists(dirName)
		isBlockDevice := info.Mode()&os.ModeDevice != 0 && info.Mode()&syscall.S_IFBLK != 0
		isCharDevice := info.Mode()&os.ModeDevice != 0 && info.Mode()&syscall.S_IFCHR != 0
		isSocket := info.Mode()&os.ModeSocket != 0
		isPipe := info.Mode()&os.ModeNamedPipe != 0
		isSetuid := info.Mode()&os.ModeSetuid != 0
		isSetgid := info.Mode()&os.ModeSetgid != 0
		isStickyDir := info.IsDir() && (info.Mode()&os.ModeSticky != 0)

		file := MyLSFiles{
			Name:          fileInfo.Name(),
			IsDir:         info.IsDir(),
			IsExec:        isExec,
			IsLink:        isLink,
			IsBroken:      isBroken,
			IsBlockDevice: isBlockDevice,
			IsCharDevice:  isCharDevice,
			IsSocket:      isSocket,
			IsPipe:        isPipe,
			IsSetuid:      isSetuid,
			IsSetgid:      isSetgid,
			IsStickyDir:   isStickyDir,
			Size:          info.Size(),
			ModTime:       info.ModTime(),
			Mode:          info.Mode(),
		}

		// If long listing is requested, format accordingly.
		if lFlag {
			fmt.Println(formatLongEntry(file, info))
		} else {
			printFiles([]MyLSFiles{file})
		}
		return
	}

	myFS := os.DirFS(dirName)
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
		IsBlockDevice := info.Mode()&os.ModeDevice != 0 && info.Mode()&syscall.S_IFBLK != 0
		IsCharDevice := info.Mode()&os.ModeDevice != 0 && info.Mode()&syscall.S_IFCHR != 0
		IsSocket := info.Mode()&os.ModeSocket != 0
		IsPipe := info.Mode()&os.ModeNamedPipe != 0
		// IsOrphanedLink := info.Mode()&os.ModeSymlink != 0
		IsSetuid := info.Mode()&os.ModeSetuid != 0
		IsSetgid := info.Mode()&os.ModeSetgid != 0
		IsStickyDir := info.Mode()&os.ModeDir != 0 && info.Mode()&os.ModeSticky != 0

		files = append(files, MyLSFiles{
			Name:          fileName,
			IsDir:         entry.IsDir(),
			IsExec:        isExec,
			IsLink:        isLink,
			IsBroken:      isBroken,
			IsBlockDevice: IsBlockDevice,
			IsCharDevice:  IsCharDevice,
			IsSocket:      IsSocket,
			IsPipe:        IsPipe,
			// IsOrphanedLink: IsOrphanedLink,
			IsSetuid:    IsSetuid,
			IsSetgid:    IsSetgid,
			IsStickyDir: IsStickyDir,
			Size:        info.Size(),
			ModTime:     info.ModTime(),
			Mode:        info.Mode(),
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
		fmt.Println(dirName + ":")
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

	}

	if !lFlag {
		printFiles(files)
	}

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
