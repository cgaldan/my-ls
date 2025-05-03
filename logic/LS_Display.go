package logic

import (
	"fmt"
	"ls/utils"
	"path/filepath"
	"strings"
	"time"
)

var punctuationMarks = []string{
	"!", "\"", "#", "$", "%", "&", "(", ")", "*", "+", ",",
	"/", ":", ";", "<", "=", ">", "?", "@", "[", "\\", "]", "^", "`", "{", "|", "}", "~", " ",
}

func formatFileNames(fileName string) string {
	if strings.Contains(fileName, "'") {
		fileName = fmt.Sprintf("\"%s\"", fileName)
	}

	for _, mark := range punctuationMarks {
		if strings.Contains(fileName, mark) {
			fileName = fmt.Sprintf("'%s'", fileName)
			break
		}
	}
	return fileName
}

// formatLongEntry returns a detailed string for a file, including extra metadata.
// It retrieves the number of links, owner, and group information from the file's syscall.Stat_t.
func formatLongEntry(file MyLSFiles, lenNLink int, maxOwner, maxGroup, maxSize int) string {
	permission := getPermission(file)

	modTime := formatTime(file.ModTime)

	size := fmt.Sprintf("%*d", maxSize, file.Size)

	fileName := formatFileNames(file.Name)

	if file.IsBlockDevice || file.IsCharDevice {
		size = fmt.Sprintf("%*d, %*d",
			len(fmt.Sprint(file.MajorNumber)), file.MajorNumber,
			len(fmt.Sprint(file.MinorNumber)), file.MinorNumber,
		)
	}

	fileColor := file.GetColor()
	if file.IsLink {
		targetColor := reset
		if file.TargetFile != nil {
			targetColor = file.TargetFile.GetColor()
		}
		fileName = fmt.Sprintf("%s%s%s -> %s%s%s", fileColor, fileName, reset, targetColor, file.LinkTarget, reset)
	}

	return fmt.Sprintf("%10s %*d %-*s %-*s %s %12s  %s%s%s",
		permission,
		lenNLink, file.NLink,
		maxOwner, file.OwnerName,
		maxGroup, file.GroupName,
		size,
		modTime,
		fileColor,
		fileName,
		reset,
	)
}

func padColoredString(uncolored, colored string, width int) string {
	// Calculate the number of spaces needed based on the visible (uncolored) length.
	padLen := max(width-len(uncolored), 0)

	return colored + strings.Repeat(" ", padLen)
}

func printFiles(files []MyLSFiles) {
	width, err := utils.GetTerminalWidth()
	if err != nil {
		width = 80
	}

	names := make([]string, len(files))
	coloredNames := make([]string, len(files))

	totalLen := 0
	var allnames []string
	for i, file := range files {
		displayName := file.Name
		displayName = formatFileNames(displayName)
		names[i] = displayName
		coloredNames[i] = file.GetColor() + displayName + reset
		names[i] = displayName

		totalLen += len(displayName)
		allnames = append(allnames, displayName)
	}

	columns := max(utils.GetColumns(len(files), width, allnames), 1)

	rows := (len(files) + columns - 1) / columns

	widthOfColumns := utils.WidthOfEachColumn(rows, columns, allnames)

	for row := range rows {
		for column := range columns {
			index := column*rows + row
			if index < len(files) {
				fmt.Print(padColoredString(names[index], coloredNames[index], widthOfColumns[column]+2))
			}
		}
	}
}

func printDirHeader(dirName string) {
	var header string

	if dirName == "." {
		header = "."
	} else if !filepath.IsAbs(dirName) && !strings.HasPrefix(dirName, "./") {
		header = "./" + dirName
	} else {
		header = dirName
	}
	fmt.Println(header + ":")
}

func getPermission(file MyLSFiles) string {
	permission := file.Mode.String()

	if file.IsLink {
		permission = strings.Replace(permission, "L", "l", 1)
	}
	if file.IsBlockDevice || file.IsCharDevice {
		permission = strings.Replace(permission, "D", "", 1)
	}
	if file.IsSetuid {
		permission = strings.Replace(permission, "u", "-", 1)
	}
	if file.IsSetgid {
		permission = strings.Replace(permission, "g", "-", 1)
	}

	return permission
}

func formatTime(modTime time.Time) string {
	now := time.Now()
	sixMonthsAgo := now.AddDate(0, -6, 0)

	if now.Sub(modTime) > 0 && modTime.After(sixMonthsAgo) {
		return modTime.Format("Jan _2 15:04")
	} else {
		return modTime.Format("Jan _2  2006")
	}
}
