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
	var fileEntries []data.MyLSFiles
	var dirPaths []string

	// Separate files and directories.
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			fmt.Printf("myls: cannot access '%s': No such file or directory\n", path)
			continue
		}
		if info.IsDir() {
			dirPaths = append(dirPaths, path)
		} else {
			entry := GetFileAttributes(path, info, true)
			fileEntries = append(fileEntries, entry)
		}
	}

	if len(fileEntries) > 0 {
		printFilesDetails(fileEntries, lFlag, rFlag, tFlag)
	}

	if len(dirPaths) > 0 && len(fileEntries) > 0 {
		fmt.Println()
	}

	// Process directories
	for i, dir := range dirPaths {
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
		dotFile := GetFileAttributes(".", fileInfo, true)
		files = append(files, dotFile)
		if stat, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
			totalBlocks += int64(stat.Blocks)
		}

		if parentInfo, err := os.Stat(dirName + "/.."); err == nil {
			parentFile := GetFileAttributes(dirName+"/..", parentInfo, false)
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

		file = GetFileAttributes(utils.Join(dirName, fileName), info, false)
		files = append(files, file)

		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			totalBlocks += int64(stat.Blocks)
		}

		if RFlag && file.IsDir {
			subDirs = append(subDirs, file)
			sortpkg.SortFiles(&subDirs, tFlag, rFlag)
		}
		UpdateMaxNlink(&maxNlink, file)
	}

	sortpkg.SortFiles(&files, tFlag, rFlag)

	if RFlag {
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
			fmt.Println() //////////////////////////////////////////////
		}
	} else {
		printFiles(files)
		if len(files) > 0 {
			fmt.Println() ////////////////////////////////////////
		}
	}

	if RFlag {
		for _, subDir := range subDirs {
			fmt.Println()
			subDirPath := utils.Join(dirName, subDir.Name)
			TheMainLS(subDirPath, lFlag, RFlag, aFlag, rFlag, tFlag)
		}
	}
}

func printFilesDetails(files []data.MyLSFiles, lFlag bool, rFlag bool, tFlag bool) {
	sortpkg.SortFiles(&files, tFlag, rFlag)
	maxOwner, maxGroup, maxsize, maxMajor, maxMinor := CalculateMaxWidth(files)
	if lFlag {
		for _, file := range files {
			fmt.Println(FormatLongEntry(file, 0, maxOwner, maxGroup, maxsize, maxMajor, maxMinor)) ////////////////////////////
		}
	} else {
		printFiles(files)
		fmt.Println() //////////////////////////////////////
	}
}
