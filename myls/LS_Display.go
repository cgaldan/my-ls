package myls

import (
	"fmt"
	"ls/utils"
	"strings"
)

func formatFileNames(fileName string) string {
	if strings.Contains(fileName, "'") { // We have to add more punctuation marks here
		fileName = fmt.Sprintf("\"%s\"", fileName)
	}
	if strings.Contains(fileName, " ") {
		fileName = fmt.Sprintf("'%s'", fileName)
	}

	return fileName
}

// formatLongEntry returns a detailed string for a file, including extra metadata.
// It retrieves the number of links, owner, and group information from the file's syscall.Stat_t.
func formatLongEntry(file MyLSFiles, lenNLink int, lenSize int) string {
	permission := file.Mode.String()

	modTim := file.ModTime.Format("Jan 2 15:04")
	modMonth := strings.Split(modTim, " ")[0]
	modMonNum := strings.Split(modTim, " ")[1]
	modTime := strings.Split(modTim, " ")[2]

	size := fmt.Sprintf("%d", file.Size)

	fileName := formatFileNames(file.Name)

	return fmt.Sprintf("%s %*d %s %s %*s %3s %2s %5s %s%s%s", permission, lenNLink, file.NLink, file.OwnerName, file.GroupName, lenSize, size, modMonth, modMonNum, modTime, file.GetColor(), fileName, reset)
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
		if strings.Contains(file.Name, " ") {
			displayName = fmt.Sprintf("'%s'", file.Name)
		}
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
		fmt.Println()
	}
}
