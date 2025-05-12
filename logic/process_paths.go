package logic

import (
	"fmt"
	"ls/data"
	"ls/sortpkg"
	"ls/utils"
	"os"
	"strings"
	"syscall"
)

func ProcessPaths(paths []string, lFlag, RFlag, aFlag, rFlag, tFlag bool) {
	var allEntries []data.MyLSFiles

	// Separate files and directories.
	for _, path := range paths {
		info, err := os.Lstat(path)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("myls: cannot access '%s': No such file or directory\n", path)
			} else {
				fmt.Printf("myls: cannot access '%s': Not a directory\n", path)
			}
			continue
		}
		entry := GetFileAttributes(path, info, true, 0)
		allEntries = append(allEntries, entry)
	}

	var files []data.MyLSFiles
	var dirs []data.MyLSFiles

	for _, entry := range allEntries {
		if entry.IsLink {
			if entry.TargetFile.IsDir && !lFlag {
				dirs = append(dirs, entry)
			} else {
				files = append(files, entry)
			}
			continue
		}
		if entry.IsDir {
			dirs = append(dirs, entry)
		} else {
			files = append(files, entry)
		}
	}

	// Sort files and directories
	sortpkg.SortFiles(&files, tFlag, rFlag)
	sortpkg.SortFiles(&dirs, tFlag, rFlag)

	if len(files) > 0 {
		printFilesDetails(files, lFlag, rFlag, tFlag)
		if len(dirs) > 0 {
			fmt.Println()
		}
	}
	// // Process directories
	for i, dir := range dirs {
		if len(allEntries) > 1 && !RFlag {
			fmt.Printf("%s:\n", dir.Name)
		}
		processDirectory(dir.Name, lFlag, RFlag, aFlag, rFlag, tFlag)
		if i != len(dirs)-1 {
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
func processDirectory(dirName string, lFlag, RFlag, aFlag, rFlag, tFlag bool) {
	var files []data.MyLSFiles
	var subDirs []data.MyLSFiles

	fileInfo, err := os.Lstat(dirName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("myls: cannot access '%s': No such file or direcory", dirName)
		} else {
			fmt.Print("Error:", err)
		}
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
	var maxNlink int

	if aFlag {
		dotFile := GetFileAttributes(".", fileInfo, true, 0)
		files = append(files, dotFile)
		if stat, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
			totalBlocks += int64(stat.Blocks)
		}

		if parentInfo, err := os.Stat(dirName + "/.."); err == nil {
			parentFile := GetFileAttributes(dirName+"/..", parentInfo, false, 0)
			files = append(files, parentFile)
			if stat, ok := parentInfo.Sys().(*syscall.Stat_t); ok {
				totalBlocks += int64(stat.Blocks)
			}
			UpdateMaxNlink(&maxNlink, parentFile)
		}
	}

	var file data.MyLSFiles
	var fileName string
	for _, entry := range entries {
		fileName = entry.Name()
		if !aFlag && strings.HasPrefix(fileName, ".") {
			continue
		}

		info, err := os.Lstat(utils.Join(dirName, fileName))
		if err != nil {
			continue
		}

		file = GetFileAttributes(utils.Join(dirName, fileName), info, false, 0)
		files = append(files, file)

		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			totalBlocks += int64(stat.Blocks)
		}

		if RFlag && file.IsDir {
			subDirs = append(subDirs, file)
		}
		UpdateMaxNlink(&maxNlink, file)
	}

	sortpkg.SortFiles(&files, tFlag, rFlag)

	if RFlag {
		sortpkg.SortFiles(&subDirs, tFlag, rFlag)
		fmt.Println(printDirHeader(dirName))
	}

	if lFlag {
		fmt.Printf("total %d\n", totalBlocks/2)
		maxOwner, maxGroup, maxsize, maxMajor, maxMinor := CalculateMaxWidth(files)
		for i, file := range files {
			fmt.Print(FormatLongEntry(file, maxNlink, maxOwner, maxGroup, maxsize, maxMajor, maxMinor))
			if i != len(files)-1 {
				fmt.Println()
			}
		}
		if len(files) > 0 {
			fmt.Println()
		}
	} else {
		printFiles(files)
		if len(files) > 0 {
			fmt.Println()
		}
	}

	if RFlag {

		for _, subDir := range subDirs {
			fmt.Println()
			for dirName[len(dirName)-1] == '/' {
				dirName = strings.TrimSuffix(dirName, "/")
			}
			dirName += "/"
			subDirPath := dirName + subDir.Name
			processDirectory(subDirPath, lFlag, RFlag, aFlag, rFlag, tFlag)
		}
	}
}

func printFilesDetails(files []data.MyLSFiles, lFlag bool, rFlag bool, tFlag bool) {
	sortpkg.SortFiles(&files, tFlag, rFlag)
	maxOwner, maxGroup, maxsize, maxMajor, maxMinor := CalculateMaxWidth(files)
	if lFlag {
		for _, file := range files {
			fmt.Println(FormatLongEntry(file, 0, maxOwner, maxGroup, maxsize, maxMajor, maxMinor))
		}
	} else {
		printFiles(files)
		fmt.Println()
	}
}
