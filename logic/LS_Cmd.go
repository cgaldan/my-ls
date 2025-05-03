package logic

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
	blue2   = "\033[34m"   // Blue
	green   = "\033[1;32m" // Bright Green for executable files
	cyan    = "\033[36m"   // Cyan for audio files
	cyan2   = "\033[1;36m" // Bright Cyan for symbolic links
	red     = "\033[1;31m" // Bright Red for broken symbolic links
	magenta = "\033[1;35m" // Bright Magenta for image files
	yellow  = "\033[1;33m" // Bright Yellow
	black   = "\033[30m"   // Black
	reset   = "\033[0m"    // Resets to the default terminal color

	// Background colors
	bgBlack  = "\033[40m"
	bgRed    = "\033[41m"
	bgYellow = "\033[43m"
	bgBlue   = "\033[44m"
	bgGreen  = "\033[42m"
)

type MyLSFiles struct {
	Name            string
	IsDir           bool
	IsExec          bool
	IsLink          bool
	LinkTarget      string
	TargetFile      *MyLSFiles
	IsBroken        bool
	IsBlockDevice   bool
	IsCharDevice    bool
	MajorNumber     uint32
	MinorNumber     uint32
	IsSocket        bool
	IsPipe          bool
	IsSetuid        bool
	IsSetgid        bool
	IsStickyDir     bool
	IsOtherWritable bool
	Size            int64
	ModTime         time.Time
	Mode            fs.FileMode
	OwnerName       string
	GroupName       string
	NLink           uint64
}

// GetColor returns the appropriate ANSI color code based on the file type.
//   - Directories are displayed in blue.
//   - Executable files are displayed in green.
//   - Symbolic links are displayed in cyan.
//   - Broken symbolic links are displayed in red.
//   - Regular files are displayed in the default terminal color.
func (file MyLSFiles) GetColor() string {
	if file.IsBroken {
		return bgBlack + red
	}
	if file.IsStickyDir && file.IsOtherWritable {
		return bgGreen + black
	}
	if file.IsOtherWritable {
		return bgGreen + blue2
	}
	if file.IsStickyDir {
		return bgBlue
	}
	if file.IsSetuid {
		return bgRed
	}
	if file.IsSetgid {
		return bgYellow + black
	}
	if file.IsLink {
		return cyan2
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

	ext := strings.ToLower(filepath.Ext(file.Name))

	switch ext {
	case // Image files
		".jpg", ".jpeg", ".gif", ".bmp", ".pbm", ".pgm", ".ppm", ".tga",
		".xbm", ".xpm", ".tif", ".tiff", ".png", ".svg", ".svgz", ".mng",
		".pcx",
		// Video files
		".mov", ".mpg", ".mpeg", ".m2v", ".mkv", ".webm", ".ogm", ".mp4", ".m4v",
		".mp4v", ".vob", ".qt", ".nuv", ".wmv", ".asf", ".rm", ".rmvb",
		".flc", ".avi", ".fli", ".flv", ".gl", ".dl", ".xcf", ".xwd",
		".yuv", ".cgm", ".emf", ".axv", ".anx", ".ogv", ".ogx":
		return magenta
	case // audio files
		".aac", ".au", ".flac", ".mid", ".midi", ".mka", ".mp3", ".mpc",
		".ogg", ".ra", ".wav", ".axa", ".oga", ".spx", ".xspf":
		return cyan
	case // compressed files
		".tar", ".tgz", ".arj", ".taz", ".lzh", ".lzma", ".tlz", ".txz",
		".zip", ".z", ".Z", ".dz", ".gz", ".lz", ".xz", ".bz2", ".bz", ".tbz", ".tbz2",
		".tz", ".deb", ".rpm", ".jar", ".rar", ".ace", ".zoo", ".cpio",
		".7z", ".rz", ".cab", ".war", ".ear", ".sar":
		return red
	}
	return reset
}

func ProcessPaths(paths []string, lFlag, RFlag, aFlag, rFlag, tFlag bool) {
	var filePaths []string
	var dirPaths []string

	// Separate files and directories.
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			// For simplicity, print error and skip that path.
			fmt.Fprintf(os.Stderr, "myls: cannot access '%s': %v\n", path, err)
			continue
		}
		if info.IsDir() {
			dirPaths = append(dirPaths, path)
		} else {
			filePaths = append(filePaths, path)
		}
	}

	// Process files: Print them together in one block.
	if len(filePaths) > 0 {
		// If you want to support sorting files similarly to how you do for directories,
		// you could convert them to a slice of your MyLSFiles and sort them.
		for i, file := range filePaths {
			fileInfo, err := os.Stat(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "myls: error reading file '%s': %v\n", file, err)
				continue
			}
			// Use your existing helper function for printing file details.
			printFileDetails(file, fileInfo, lFlag)
			if i == len(filePaths)-1 && !lFlag {
				fmt.Println()
			}
		}
		// If there are directories also, separate the file block from directories with a blank line.
		if len(dirPaths) > 0 {
			fmt.Println()
		}
	}

	// Process directories
	for i, dir := range dirPaths {
		// If more than one path is provided, print a header to mimic ls behavior.
		// (You can also check that there's more than one directory.)
		if len(paths) > 1 && !RFlag {
			fmt.Printf("%s:\n", dir)
		}
		TheMainLS(dir, lFlag, RFlag, aFlag, rFlag, tFlag)
		if i != len(dirPaths)-1 {
			fmt.Println()
		}
	}
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

	fileInfo, err := os.Lstat(dirName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("myls: cannot access '%s': No such file or direcory", dirName)
		} else {
			fmt.Print("Error:", err)
		}
		return
	}

	if !fileInfo.IsDir() {
		printFileDetails(dirName, fileInfo, lFlag)
		return
	}

	entries, err := os.ReadDir(dirName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("myls: cannot access '%s': No such file or direcory", dirName)
		} else {
			fmt.Print("Error:", err)
		}
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

		if parentInfo, err := os.Stat(dirName + "/.."); err == nil {
			parentFile := getFileAttributes(dirName+"/..", parentInfo)
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
		printDirHeader(dirName)
	}

	if lFlag {
		fmt.Printf("total %d\n", totalBlocks/2)
		maxOwner, maxGroup, maxSize := calculateMaxWidth(files)
		for i, file := range files {
			fmt.Print(formatLongEntry(file, maxNlink, maxOwner, maxGroup, maxSize))
			if i != len(files)-1 {
				fmt.Println()
			}
		}
		fmt.Println() //////////////////////////////////////////////
	} else {
		printFiles(files)
		if len(files) > 0 {
			fmt.Println() ////////////////////////////////////////
		}
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
		fmt.Println(formatLongEntry(file, len(strconv.Itoa(int(file.NLink))), len(strconv.Itoa(int(file.Size))), len(file.OwnerName), len(file.GroupName))) ////////////////////////////
	} else {
		printFiles([]MyLSFiles{file})
		// fmt.Println() //////////////////////////////////////
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

func calculateMaxWidth(files []MyLSFiles) (maxOwner, maxGroup, maxSize int) {
	for _, file := range files {
		if len(file.OwnerName) > maxOwner {
			maxOwner = len(file.OwnerName)
		}
		if len(file.GroupName) > maxGroup {
			maxGroup = len(file.GroupName)
		}

		var sizeField string
		if file.IsBlockDevice || file.IsCharDevice {
			sizeField = fmt.Sprintf("%3d, %3d", file.MajorNumber, file.MinorNumber)
		} else {
			sizeField = fmt.Sprintf("%d", file.Size)
		}

		if len(sizeField) > maxSize {
			maxSize = len(sizeField)
		}
	}

	return maxOwner, maxGroup, maxSize
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

	var targetFile *MyLSFiles
	var targetPath string
	var err error
	if info.Mode()&os.ModeSymlink != 0 {
		targetPath, err = os.Readlink(path)
		if err == nil {
			absTarget := filepath.Join(filepath.Dir(path), targetPath)

			if targetInfo, err := os.Lstat(absTarget); err == nil {

				tf := getFileAttributes(absTarget, targetInfo)
				targetFile = &tf
			}
		}
	}

	var major, minor uint32
	if stat != nil {
		major = uint32((stat.Rdev >> 8) & 0xFF) // Linux/Unix specific
		minor = uint32(stat.Rdev & 0xFF)
	}

	return MyLSFiles{
		Name:            filepath.Base(path),
		IsDir:           info.IsDir(),
		IsExec:          !info.IsDir() && (info.Mode().Perm()&0o111 != 0),
		IsLink:          info.Mode()&os.ModeSymlink != 0,
		LinkTarget:      targetPath,
		TargetFile:      targetFile,
		IsBroken:        info.Mode()&os.ModeSymlink != 0 && !exists(path),
		IsBlockDevice:   info.Mode()&os.ModeType == os.ModeDevice,
		IsCharDevice:    info.Mode()&os.ModeType == (os.ModeDevice | os.ModeCharDevice),
		MajorNumber:     major,
		MinorNumber:     minor,
		IsSocket:        info.Mode()&os.ModeSocket != 0,
		IsPipe:          info.Mode()&os.ModeNamedPipe != 0,
		IsSetuid:        info.Mode()&os.ModeSetuid != 0,
		IsSetgid:        info.Mode()&os.ModeSetgid != 0,
		IsStickyDir:     info.IsDir() && info.Mode()&os.ModeSticky != 0,
		IsOtherWritable: info.IsDir() && info.Mode()&0o002 != 0,
		Size:            info.Size(),
		ModTime:         info.ModTime(),
		Mode:            info.Mode(),
		OwnerName:       ownerName,
		GroupName:       groupName,
		NLink:           nlink,
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
