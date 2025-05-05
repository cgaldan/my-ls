package logic

import (
	"fmt"
	"ls/data"
	"ls/utils"
	"strings"
)

const reset = "\033[0m"

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

func padColoredString(uncolored, colored string, width int) string {
	// Calculate the number of spaces needed based on the visible (uncolored) length.
	padLen := max(width-len(uncolored), 0)

	return colored + strings.Repeat(" ", padLen)
}

func printFiles(files []data.MyLSFiles) {
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
		if rows > 1 && row != rows-1 {
			fmt.Println()
		}
	}
}

func printDirHeader(dirName string) string {
	var header string

	if dirName == "." {
		header = "."
	} else {
		header = dirName
	}
	return header + ":"
}